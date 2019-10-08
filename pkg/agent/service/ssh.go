package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"
	"unsafe"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/pkg/errors"
	gossh "golang.org/x/crypto/ssh"
)

const (
	authorizedKeysFilename = "authorized_keys"

	// Simple script to start a preferred shell
	// On Debian and Ubuntu /bin/sh links to dash, whereas bash is what's actually preferred
	// This is fairly hacky and there's likely a better approach to determining the preferred shell
	entrypoint = `if [ "$(readlink /bin/sh)" = "dash" ] && [ -f "/bin/bash" ]; then exec /bin/bash; else exec /bin/sh; fi`
)

func (s *Service) ssh(w http.ResponseWriter, r *http.Request) {
	if s.variables.GetDisableSSH() {
		http.Error(w, "SSH is disabled", http.StatusForbidden)
		return
	}

	conn := conncontext.GetConn(r)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if s.variables.GetDisableSSH() {
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	signer, err := generateSigner()
	if err != nil {
		http.Error(w, errors.Wrap(err, "generate signer").Error(), http.StatusInternalServerError)
		return
	}

	// TODO: refactor this
	// The variables package should be the only one reading variables
	authorizedKeysLocation := path.Join(s.confDir, authorizedKeysFilename)

	var options []ssh.Option
	if _, err := os.Stat(authorizedKeysLocation); err == nil {
		authorizedKeyBytes, err := ioutil.ReadFile(authorizedKeysLocation)
		if err != nil {
			log.WithError(err).Error("read authorized keys")
			return
		}
		authorizedKey, _, _, _, err := gossh.ParseAuthorizedKey(authorizedKeyBytes)
		if err != nil {
			log.WithError(err).Error("parse authorized keys")
			return
		}
		options = []ssh.Option{
			ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
				return ssh.KeysEqual(key, authorizedKey)
			}),
		}
	} else if !os.IsNotExist(err) {
		log.WithError(err).Error("check for authorized keys")
		return
	}

	sshServer := &ssh.Server{
		Handler:         handler(ctx),
		RequestHandlers: ssh.DefaultRequestHandlers,
		ChannelHandlers: ssh.DefaultChannelHandlers,
		HostSigners:     []ssh.Signer{signer},
	}

	for _, option := range options {
		if err = sshServer.SetOption(option); err != nil {
			http.Error(w, errors.Wrap(err, "set SSH option").Error(), http.StatusInternalServerError)
		}
	}

	sshServer.HandleConn(conn)
}

func handler(ctx context.Context) func(s ssh.Session) {
	return func(s ssh.Session) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var cmd *exec.Cmd
		if _, err := os.Stat("/usr/bin/nsenter"); os.IsNotExist(err) {
			cmd = exec.CommandContext(ctx, "/bin/sh", "-c", entrypoint)
		} else {
			cmd = exec.CommandContext(ctx, "/usr/bin/nsenter", "-t", "1", "-a", "/bin/sh", "-c", entrypoint)
		}

		ptyReq, winCh, isPty := s.Pty()
		if !isPty {
			return
		}

		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))

		f, err := pty.Start(cmd)
		if err != nil {
			log.WithError(err).Error("start PTY")
			return
		}

		go func() {
			for win := range winCh {
				syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
					uintptr(unsafe.Pointer(&struct {
						h, w, x, y uint16
					}{
						uint16(win.Height), uint16(win.Width), 0, 0,
					})))
			}
		}()

		go io.Copy(f, s)
		io.Copy(s, f)
	}
}

func generateSigner() (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return gossh.NewSignerFromKey(key)
}

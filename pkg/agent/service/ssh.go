package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/pkg/errors"
)

const (
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

	signer, err := s.getSigner()
	if err != nil {
		http.Error(w, errors.Wrap(err, "generate signer").Error(), http.StatusInternalServerError)
		return
	}

	sshServer := &ssh.Server{
		Handler:         sshServerHandler(ctx),
		RequestHandlers: ssh.DefaultRequestHandlers,
		ChannelHandlers: ssh.DefaultChannelHandlers,
		HostSigners:     []ssh.Signer{signer},
	}

	var options []ssh.Option
	if len(s.variables.GetAuthorizedSSHKeys()) > 0 {
		options = []ssh.Option{
			ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
				for _, authorizedKey := range s.variables.GetAuthorizedSSHKeys() {
					if ssh.KeysEqual(key, authorizedKey) {
						return true
					}
				}
				return false
			}),
		}
	}

	for _, option := range options {
		if err = sshServer.SetOption(option); err != nil {
			http.Error(w, errors.Wrap(err, "set SSH option").Error(), http.StatusInternalServerError)
		}
	}

	sshServer.HandleConn(conn)
}

func sshServerHandler(ctx context.Context) func(s ssh.Session) {
	return func(s ssh.Session) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		innerCommand := s.RawCommand()
		if innerCommand == "" {
			innerCommand = entrypoint
		}

		command := []string{"/bin/sh", "-c", innerCommand}

		cmd := exec.CommandContext(ctx, command[0], command[1:]...)

		ptyReq, winCh, isPty := s.Pty()
		if isPty {
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
		} else {
			cmd.Stdout = s
			cmd.Stderr = s
			if err := cmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					s.Exit(exitError.ExitCode())
					return
				}
				log.WithError(err).Error("run SSH command")
				return
			}
		}
	}
}

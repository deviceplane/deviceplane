package connector

import (
	"context"
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
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/gliderlabs/ssh"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	gossh "golang.org/x/crypto/ssh"
)

const (
	authorizedKeysFilename = "authorized_keys"

	// Simple script to start a preferred shell
	// On Debian and Ubuntu /bin/sh links to dash, whereas bash is what's actually preferred
	// This is fairly hacky and there's likely a better approach to determining the preffered shell
	entrypoint = `if [ "$(readlink /bin/sh)" = "dash" ] && [ -f "/bin/bash" ]; then exec /bin/bash; else exec /bin/sh; fi`
)

type Connector struct {
	client    *agent_client.Client // TODO: interface
	variables variables.Interface
	confDir   string
}

func NewConnector(
	client *agent_client.Client, variables variables.Interface,
	confDir string,
) *Connector {
	return &Connector{
		client:    client,
		variables: variables,
		confDir:   confDir,
	}
}

func (c *Connector) Do() {
	if c.variables.GetDisableSSH() {
		return
	}

	conn, err := c.client.InitiateDeviceConnection(context.TODO())
	if err != nil {
		log.WithError(err).Error("initiate connection")
		return
	}

	listener := revdial.NewListener(conn, func(ctx context.Context, path string) (*websocket.Conn, *http.Response, error) {
		return c.client.Revdial(ctx, path)
	})
	defer listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if c.variables.GetDisableSSH() {
					listener.Close()
				}
				if listener.Closed() {
					cancel()
					return
				}
			}
		}
	}()

	ssh.Handle(func(s ssh.Session) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var cmd *exec.Cmd
		if _, err = os.Stat("/usr/bin/nsenter"); os.IsNotExist(err) {
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
	})

	authorizedKeysLocation := path.Join(c.confDir, authorizedKeysFilename)

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

	if err = ssh.Serve(listener, nil, options...); err != nil {
		log.WithError(err).Error("serve SSH")
	}
}

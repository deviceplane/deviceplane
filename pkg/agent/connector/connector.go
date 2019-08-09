package connector

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/apex/log"
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
)

type Connector struct {
	client    *agent_client.Client // TODO: interface
	variables variables.Interface
}

func NewConnector(client *agent_client.Client, variables variables.Interface) *Connector {
	return &Connector{
		client:    client,
		variables: variables,
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

	listener := revdial.NewListener(conn, func(ctx context.Context) (net.Conn, error) {
		return c.client.Dial(ctx)
	})

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
			cmd = exec.CommandContext(ctx, "/bin/sh")
		} else {
			cmd = exec.CommandContext(ctx, "/usr/bin/nsenter", "-t", "1", "-m", "-u", "-i", "-n", "-p")
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

	if err = ssh.Serve(listener, nil); err != nil {
		log.WithError(err).Error("serve SSH")
	}
}

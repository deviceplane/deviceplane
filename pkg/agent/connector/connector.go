package connector

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/apex/log"
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
)

type Connector struct {
	client *agent_client.Client // TODO: interface
}

func NewConnector(client *agent_client.Client) *Connector {
	return &Connector{
		client: client,
	}
}

func (c *Connector) Do() {
	conn, err := c.client.InitiateDeviceConnection(context.TODO())
	if err != nil {
		log.WithError(err).Error("initiate connection")
		return
	}

	listener := revdial.NewListener(conn, func(ctx context.Context) (net.Conn, error) {
		return c.client.Dial(ctx)
	})

	ssh.Handle(func(s ssh.Session) {
		cmd := exec.Command("nsenter", "-t", "1", "-m", "-u", "-i", "-n", "-p")

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
				setWinsize(f, win.Width, win.Height)
			}
		}()

		go io.Copy(f, s)
		io.Copy(s, f)
	})

	if err = ssh.Serve(listener, nil); err != nil {
		log.WithError(err).Error("SSH server")
	}
}

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

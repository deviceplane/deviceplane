package netns

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/vishvananda/netns"
)

const (
	timeout = time.Second
)

type request struct {
	ctx         context.Context
	containerID string
	port        int
	path        string
}

type response struct {
	response *http.Response
	err      error
}

type Manager struct {
	engine engine.Engine
	in     chan request
	out    chan response
}

func NewManager(engine engine.Engine) *Manager {
	return &Manager{
		engine: engine,
		in:     make(chan request),
		out:    make(chan response),
	}
}

func (m *Manager) Start() {
	go func() {
		runtime.LockOSThread()
		for {
			select {
			case request := <-m.in:
				ctx, cancel := context.WithTimeout(request.ctx, timeout)
				m.out <- m.processRequest(ctx, request)
				cancel()
			}
		}
	}()
}

func (m *Manager) ProcessRequest(
	ctx context.Context, containerID string, port int, path string,
) (*http.Response, error) {
	m.in <- request{
		ctx:         ctx,
		containerID: containerID,
		port:        port,
		path:        path,
	}
	resp := <-m.out
	return resp.response, resp.err
}

func (m *Manager) processRequest(ctx context.Context, req request) response {
	inspectResponse, err := m.engine.InspectContainer(ctx, req.containerID)
	if err != nil {
		return response{
			err: err,
		}
	}

	containerNamespace, err := netns.GetFromPid(inspectResponse.PID)
	if err != nil {
		return response{
			err: err,
		}
	}
	defer containerNamespace.Close()

	if err := netns.Set(containerNamespace); err != nil {
		return response{
			err: err,
		}
	}

	var dialer net.Dialer
	conn, err := dialer.DialContext(
		ctx, "tcp", fmt.Sprintf("127.0.0.1:%d", req.port),
	)
	if err != nil {
		return response{
			err: err,
		}
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx, "GET", string(req.path), nil,
	)
	if err != nil {
		return response{
			err: err,
		}
	}

	if err := httpRequest.Write(conn); err != nil {
		return response{
			err: err,
		}
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), httpRequest)
	if err != nil {
		return response{
			err: err,
		}
	}

	return response{
		response: resp,
	}
}

package connman

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/deviceplane/deviceplane/pkg/revdial"
)

var (
	ErrNoConnection = errors.New("no connection")
)

type ConnectionManager struct {
	deviceConnections map[string]net.Conn
	deviceDialers     map[string]*revdial.Dialer
	lock              sync.RWMutex
}

func New() *ConnectionManager {
	return &ConnectionManager{
		deviceConnections: make(map[string]net.Conn),
		deviceDialers:     make(map[string]*revdial.Dialer),
	}
}

func (m *ConnectionManager) Set(key string, conn net.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.deviceConnections[key] = conn
	m.deviceDialers[key] = revdial.NewDialer(conn, "/revdial")
}

func (m *ConnectionManager) Join(key string, conn net.Conn) error {
	m.lock.RLock()
	dialer, ok := m.deviceDialers[key]
	if !ok {
		m.lock.RUnlock()
		return ErrNoConnection
	}
	m.lock.RUnlock()

	otherConn, err := dialer.Dial(context.Background())
	if err != nil {
		return nil
	}

	go io.Copy(otherConn, conn)
	io.Copy(conn, otherConn)

	return nil
}

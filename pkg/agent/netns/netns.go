package netns

import (
	"context"
	"runtime"

	"github.com/vishvananda/netns"

	"github.com/deviceplane/deviceplane/pkg/engine"
)

type Manager struct {
	engine engine.Engine
}

func NewManager(engine engine.Engine) *Manager {
	return &Manager{
		engine: engine,
	}
}

func (m *Manager) RunInContainerNamespace(ctx context.Context, containerID string, f func()) error {
	inspectResponse, err := m.engine.InspectContainer(ctx, containerID)
	if err != nil {
		return err
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	originalNamespace, err := netns.Get()
	if err != nil {
		return err
	}
	defer originalNamespace.Close()
	defer netns.Set(originalNamespace)

	containerNamespace, err := netns.GetFromPid(inspectResponse.PID)
	if err != nil {
		return err
	}
	defer containerNamespace.Close()

	if err := netns.Set(containerNamespace); err != nil {
		return err
	}

	f()

	return nil
}

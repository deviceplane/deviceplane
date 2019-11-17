package supervisor

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"

	agent_utils "github.com/deviceplane/deviceplane/pkg/agent/utils"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

type PullEvent struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

type imagePuller struct {
	applicationID string
	serviceName   string
	engine        engine.Engine

	currentlyPulling atomic.Value
	progress         map[string]PullEvent
	lock             sync.RWMutex
}

func newImagePuller(
	applicationID string,
	serviceName string,
	engine engine.Engine,
) *imagePuller {
	p := &imagePuller{
		applicationID: applicationID,
		serviceName:   serviceName,
		engine:        engine,

		progress: make(map[string]PullEvent),
	}
	p.currentlyPulling.Store(false)
	return p
}

func (p *imagePuller) Pull(ctx context.Context, image string) {
	p.currentlyPulling.Store(true)
	defer p.currentlyPulling.Store(false)

	p.lock.Lock()
	p.progress = make(map[string]PullEvent)
	p.lock.Unlock()

	r, w := io.Pipe()
	go func() {
		decoder := json.NewDecoder(r)
		for {
			var event PullEvent
			if err := decoder.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			p.lock.Lock()
			p.progress[event.ID] = event
			p.lock.Unlock()
		}
	}()

	agent_utils.ImagePull(ctx, p.engine, image, w)
}

func (p *imagePuller) Progress() (map[string]PullEvent, bool) {
	if !p.currentlyPulling.Load().(bool) {
		return nil, false
	}

	p.lock.RLock()
	defer p.lock.RUnlock()

	var progressCopy map[string]PullEvent
	utils.JSONConvert(p.progress, &progressCopy)
	return progressCopy, true
}

package runner

import (
	"context"
	"sync"
	"time"
)

const (
	repeatInverval = time.Minute
)

// Manager is responsible for running all background runners
// This will play a more important role down the road when we can run
// multiple controllers and need distributed locking
type Manager struct {
	runners []Runner
}

func NewManager(runners []Runner) *Manager {
	return &Manager{
		runners: runners,
	}
}

func (m *Manager) Start() {
	go func() {
		ticker := time.NewTicker(repeatInverval)
		defer ticker.Stop()

		for {
			ctx, cancel := context.WithTimeout(context.Background(), repeatInverval/2)

			var wg sync.WaitGroup
			for _, runner := range m.runners {
				wg.Add(1)
				go func(runner Runner) {
					runner.Do(ctx)
					wg.Done()
				}(runner)
			}

			wg.Wait()
			cancel()

			select {
			case <-ticker.C:
				continue
			}
		}
	}()
}

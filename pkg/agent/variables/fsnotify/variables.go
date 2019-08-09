package fsnotify

import (
	"os"
	"path"
	"sync"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/fsnotify/fsnotify"
)

type Variables struct {
	dir           string
	lock          sync.RWMutex
	getDisableSSH bool
}

func NewVariables(dir string) *Variables {
	return &Variables{
		dir:           dir,
		getDisableSSH: false,
	}
}

func (v *Variables) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				if err = v.refresh(); err != nil {
					log.WithError(err).Error("variables refresh")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.WithError(err).Error("variables watcher error")
			}
		}
	}()

	return watcher.Add(v.dir)
}

func (v *Variables) refresh() error {
	_, err := os.Stat(path.Join(v.dir, variables.DisableSSH))

	v.lock.Lock()
	defer v.lock.Unlock()

	if err == nil {
		v.getDisableSSH = true
	} else if os.IsNotExist(err) {
		v.getDisableSSH = false
	} else {
		return err
	}

	return nil
}

func (v *Variables) GetDisableSSH() bool {
	return v.getDisableSSH
}

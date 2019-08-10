package fsnotify

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/fsnotify/fsnotify"
)

type Variables struct {
	dir           string
	lock          sync.RWMutex
	getDisableSSH *bool
}

func NewVariables(dir string) *Variables {
	return &Variables{
		dir: dir,
	}
}

func (v *Variables) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err = v.refresh(); err != nil {
		log.WithError(err).Error("variables refresh")
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
		v.getDisableSSH = &[]bool{true}[0]
	} else if os.IsNotExist(err) {
		v.getDisableSSH = &[]bool{false}[0]
	} else {
		return err
	}

	return nil
}

func (v *Variables) GetDisableSSH() bool {
	return v.getField(func() *bool {
		return v.getDisableSSH
	})
}

func (v *Variables) getField(getField func() *bool) bool {
	v.lock.RLock()
	field := getField()
	v.lock.RUnlock()
	if field != nil {
		return *field
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			v.lock.RLock()
			field := getField()
			v.lock.RUnlock()
			if field != nil {
				return *field
			}
		}
	}
}

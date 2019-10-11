package fsnotify

import (
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/crypto/ssh"
)

type Variables struct {
	dir                  string
	lock                 sync.RWMutex
	disableSSH           bool
	disableSSHSet        bool
	authorizedSSHKeys    []ssh.PublicKey
	authorizedSSHKeysSet bool
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

	v.refresh()

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				v.refresh()
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

func (v *Variables) refresh() {
	for _, refresher := range []func() error{
		v.refreshDisableSSH,
		v.refreshAuthorizedSSHKeys,
	} {
		if err := refresher(); err != nil {
			log.WithError(err).Error("variables refresh")
		}
	}
}

func (v *Variables) refreshDisableSSH() error {
	_, err := os.Stat(path.Join(v.dir, variables.DisableSSH))

	v.lock.Lock()
	defer v.lock.Unlock()

	if err == nil {
		v.disableSSH = true
		v.disableSSHSet = true
	} else if os.IsNotExist(err) {
		v.disableSSH = false
		v.disableSSHSet = true
	} else {
		return err
	}

	return nil
}

func (v *Variables) refreshAuthorizedSSHKeys() error {
	bytes, err := ioutil.ReadFile(path.Join(v.dir, variables.AuthorizedSSHKeys))

	v.lock.Lock()
	defer v.lock.Unlock()

	if err == nil {
		authorizedSSHKeys, err := parseAuthorizedKeysFile(bytes)
		if err != nil {
			return err
		}
		v.authorizedSSHKeys = authorizedSSHKeys
		v.authorizedSSHKeysSet = true
	} else if os.IsNotExist(err) {
		v.authorizedSSHKeys = make([]ssh.PublicKey, 0)
		v.authorizedSSHKeysSet = true
	} else {
		return err
	}

	return nil
}

func (v *Variables) GetDisableSSH() bool {
	v.waitFor(func() bool {
		return v.disableSSHSet
	})
	return v.disableSSH
}

func (v *Variables) GetAuthorizedSSHKeys() []ssh.PublicKey {
	v.waitFor(func() bool {
		return v.authorizedSSHKeysSet
	})
	return v.authorizedSSHKeys
}

func (v *Variables) waitFor(getField func() bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			v.lock.RLock()
			field := getField()
			v.lock.RUnlock()
			if field {
				return
			}
		}
	}
}

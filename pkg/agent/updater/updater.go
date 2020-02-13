package updater

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"
	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
	dphttp "github.com/deviceplane/deviceplane/pkg/http"
)

const (
	location        = "https://downloads.deviceplane.com/agent/%s/linux/%s/deviceplane-agent"
	downloadTimeout = time.Hour
)

type Updater struct {
	projectID  string
	version    string
	binaryPath string

	desiredVersion string
	once           sync.Once
	lock           sync.RWMutex
}

func NewUpdater(projectID, version, binaryPath string) *Updater {
	return &Updater{
		projectID:  projectID,
		version:    version,
		binaryPath: binaryPath,
	}
}

func (u *Updater) SetDesiredVersion(desiredVersion string) {
	u.lock.Lock()
	u.desiredVersion = desiredVersion
	u.lock.Unlock()

	u.once.Do(func() {
		go u.updater()
	})
}

func (u *Updater) updater() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		u.lock.RLock()
		desiredVersion := u.desiredVersion
		u.lock.RUnlock()

		if desiredVersion != "" && desiredVersion != u.version {
			ctx, cancel := dpcontext.New(context.Background(), downloadTimeout)
			defer cancel()

			if err := u.update(ctx, desiredVersion); err != nil {
				log.WithError(err).Error("update agent")
				goto cont
			}
		}

	cont:
		<-ticker.C
	}
}

func (u *Updater) update(ctx *dpcontext.Context, desiredVersion string) error {
	resp, err := dphttp.Get(ctx, fmt.Sprintf(location, desiredVersion, runtime.GOARCH))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	for _, action := range []func() error{
		func() error {
			_, err := io.Copy(f, resp.Body)
			return err
		},
		func() error {
			return f.Close()
		},
		func() error {
			return os.Chmod(f.Name(), 0755)
		},
		func() error {
			return syscall.Unlink(u.binaryPath)
		},
		func() error {
			return os.Rename(f.Name(), u.binaryPath)
		},
	} {
		if err = action(); err != nil {
			return err
		}
	}

	os.Exit(0)
	return nil
}

package service

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/apex/log"
)

func (s *Service) reboot(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(1000)
		err := exec.Command("/sbin/reboot").Run()
		if err != nil {
			log.WithError(err).Error("failed to reboot")
		}
	}()
}

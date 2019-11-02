package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func (s *Service) execute(w http.ResponseWriter, r *http.Request) {
	commandBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), codes.StatusInternalDeviceError)
		return
	}

	command, err := nsenterCommandWrapper(string(commandBytes))
	if err != nil {
		http.Error(w, err.Error(), codes.StatusInternalDeviceError)
		return
	}

	if _, ok := r.URL.Query()["background"]; ok {
		if err := exec.CommandContext(
			context.Background(), command[0], command[1:]...,
		).Start(); err != nil {
			http.Error(w, err.Error(), codes.StatusInternalDeviceError)
			return
		}
	} else {
		exitCode := 0

		if err := exec.CommandContext(
			r.Context(), command[0], command[1:]...,
		).Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
				}
			} else {
				http.Error(w, err.Error(), codes.StatusInternalDeviceError)
				return
			}
		}

		utils.Respond(w, models.ExecuteResponse{
			ExitCode: exitCode,
		})
	}
}

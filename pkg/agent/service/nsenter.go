package service

import (
	"os"
)

// nsenterCommandWrapper wraps a command in an nsenter command if the
// nsenter command is found
//
// The agent container should always contain nsenter and so a command
// getting wrapped in nsenter is the expected flow for most cases
//
// The other path is primarily there to make SSH work for local development
func nsenterCommandWrapper(command ...string) ([]string, error) {
	wrappedCommand := []string{"/bin/sh", "-c"}

	if _, err := os.Stat("/usr/bin/nsenter"); err == nil {
		wrappedCommand = append([]string{
			"/usr/bin/nsenter",
			"-t", "1",
			"-a",
		}, wrappedCommand...)
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	wrappedCommand = append(wrappedCommand, command...)

	return wrappedCommand, nil
}

package cliutils

import (
	"time"

	"github.com/hako/durafmt"
)

func DurafmtSince(d time.Time) *durafmt.Durafmt {
	diff := time.Now().Sub(d).Truncate(time.Millisecond)
	duration := durafmt.Parse(diff).LimitFirstN(1)
	return duration
}

func GetSSHArgs(args []string) (preSSH []string, postSSH []string) {
	var i int
	var hasSSH bool
	for i = 0; i < len(args); i++ {
		if i > 0 && args[i-1] == "ssh" { // Split like so: deviceplane [...] ssh [device] [post-ssh]
			hasSSH = true
			break
		}
	}

	if !hasSSH {
		return args, nil
	}

	preSSH = args[0 : i+1]
	if len(args) > i+1 {
		postSSH = args[i+1:]
	}
	return
}

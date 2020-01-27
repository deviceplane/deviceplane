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

package filtering

import (
	"regexp"

	"github.com/deviceplane/deviceplane/pkg/utils"
)

var nodePrefixRegex = regexp.MustCompile(`(?m)(^|^# HELP |^# TYPE )(node_)`)

func FilterNodePrefix(rawHostMetrics string) string {
	return utils.ReplaceAllStringSubmatchFunc(nodePrefixRegex, rawHostMetrics, func(s []string) string {
		return s[1]
	})
}

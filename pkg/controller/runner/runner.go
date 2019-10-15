package runner

import "context"

type Runner interface {
	Do(context.Context)
}

package context

import (
	"context"
	"time"
)

var _ context.Context = &Context{}

type Context struct {
	context.Context
}

func New(ctx context.Context, timeout time.Duration) (*Context, func()) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	return &Context{
		Context: ctx,
	}, cancel
}

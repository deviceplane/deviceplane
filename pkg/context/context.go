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

/*func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c *Context) Err() error {
	return c.Context.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}*/

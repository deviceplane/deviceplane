package utils

import (
	"context"
	"time"
)

func Retry(ctx context.Context, f func(context.Context) error, timeout time.Duration) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		ctx2, cancel := context.WithTimeout(ctx, timeout)
		if err := f(ctx2); err == nil {
			cancel()
			return nil
		}
		cancel()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}

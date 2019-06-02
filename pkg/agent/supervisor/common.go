package supervisor

import (
	"context"
	"time"
)

func retry(ctx context.Context, f func(context.Context) error, timeout time.Duration) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		ctx2, cancel := context.WithTimeout(ctx, timeout)
		if err := f(ctx2); err == nil {
			cancel()
			return
		}
		cancel()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			continue
		}
	}
}

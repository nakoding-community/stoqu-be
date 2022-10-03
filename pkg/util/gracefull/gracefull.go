package gracefull

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ProcessStarter func() error

type ProcessStopper func(ctx context.Context) error

func StartProcessAtBackground(ps ...func() error) {
	for _, p := range ps {
		if p != nil {
			go func(_p func() error) {
				_ = _p()
			}(p)
		}
	}
}

func StopProcessAtBackground(duration time.Duration, ps ...ProcessStopper) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	for _, p := range ps {
		if p == nil {
			continue
		}
		ctx, stop := context.WithTimeout(context.Background(), duration)
		defer stop()
		_ = p(ctx)

	}
}

func NilStopper(ctx context.Context) error { return nil }

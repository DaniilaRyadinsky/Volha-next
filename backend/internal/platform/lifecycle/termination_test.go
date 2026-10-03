package lifecycle

import (
	"context"
	"errors"
	"testing"
)

func TestWaitTermination(t *testing.T) {
	t.Parallel()

	t.Run("signal", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		got := Wait(ctx, nil, nil)

		if got.Source != "signal" || !errors.Is(got.Err, context.Canceled) {
			t.Fatalf("termination = %+v, want signal canceled", got)
		}
	})

	t.Run("server error", func(t *testing.T) {
		want := errors.New("server failed")
		serverErr := make(chan error, 1)
		serverErr <- want

		got := Wait(context.Background(), serverErr, nil)

		if got.Source != "server" || !errors.Is(got.Err, want) {
			t.Fatalf("termination = %+v, want server error", got)
		}
	})

	t.Run("metrics error", func(t *testing.T) {
		want := errors.New("metrics failed")
		metricsErr := make(chan error, 1)
		metricsErr <- want

		got := Wait(context.Background(), nil, metricsErr)

		if got.Source != "metrics" || !errors.Is(got.Err, want) {
			t.Fatalf("termination = %+v, want metrics error", got)
		}
	})

	t.Run("closed metrics channel", func(t *testing.T) {
		metricsErr := make(chan error)
		close(metricsErr)

		got := Wait(context.Background(), nil, metricsErr)

		if got.Source != "metrics" || got.Err != nil {
			t.Fatalf("termination = %+v, want closed metrics without error", got)
		}
	})
}

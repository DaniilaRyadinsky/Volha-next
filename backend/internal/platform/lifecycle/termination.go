package lifecycle

import "context"

type Termination struct {
	Source string
	Err    error
}

func Wait(ctx context.Context, serverErr <-chan error, metricsErr <-chan error) Termination {
	select {
	case err, ok := <-serverErr:
		if !ok {
			return Termination{Source: "server"}
		}
		return Termination{Source: "server", Err: err}
	case err, ok := <-metricsErr:
		if !ok {
			return Termination{Source: "metrics"}
		}
		return Termination{Source: "metrics", Err: err}
	case <-ctx.Done():
		return Termination{Source: "signal", Err: ctx.Err()}
	}
}

package publiccloud

import (
	"context"
	"fmt"
	"time"
)

// PollFunc returns the current lifecycle status string and any error
// encountered while fetching it.
type PollFunc func() (string, error)

// WaitForStatusOk polls fetch until it returns StatusOk, an error, or the
// context / timeout fires. The retry interval grows linearly up to 30s to
// keep early iterations responsive while avoiding hammering the API on long
// waits.
func WaitForStatusOk(ctx context.Context, fetch PollFunc, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	interval := 2 * time.Second

	for {
		status, err := fetch()
		if err != nil {
			return err
		}
		switch status {
		case StatusOk:
			return nil
		case StatusError:
			return fmt.Errorf("resource entered error state")
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for status=ok (last status=%q)", status)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}

		if interval < 30*time.Second {
			interval += 2 * time.Second
		}
	}
}

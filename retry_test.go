package retry

import (
	"context"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	t.Run("NoRetry", func(t *testing.T) {
		var attempt int
		fn := func() error {
			attempt++
			return nil
		}
		err := Do(fn)
		checkResult(t, nil, err, 1, attempt)
	})
	t.Run("ErrorNotRetriable", func(t *testing.T) {
		var attempt int
		fn := func() error {
			attempt++
			return context.Canceled
		}
		err := Do(fn, WithMax(3), WithIsRetriable(IsDeadlineExceededError))
		checkResult(t, context.Canceled, err, 1, attempt)
	})
	t.Run("ErrorRetriable", func(t *testing.T) {
		var attempt int
		fn := func() error {
			attempt++
			return context.DeadlineExceeded
		}
		err := Do(fn, WithMax(3), WithIsRetriable(IsDeadlineExceededError))
		checkResult(t, context.DeadlineExceeded, err, 3, attempt)
	})
	t.Run("ErrorRetriable2", func(t *testing.T) {
		var attempt int
		fn := func() error {
			attempt++
			if attempt > 1 {
				return nil
			}
			return context.DeadlineExceeded
		}
		err := Do(fn, WithMax(3), WithIsRetriable(IsDeadlineExceededError))
		checkResult(t, nil, err, 2, attempt)
	})

	t.Run("TimeoutError", func(t *testing.T) {
		var attempt int
		fn := func(ctx context.Context) error {
			attempt++
			timer := time.NewTimer(time.Millisecond * 2)
			defer timer.Stop()
			select {
			case <-timer.C:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		err := DoWithContext(context.Background(), fn,
			WithMax(3),
			WithIsRetriable(IsDeadlineExceededError),
			WithPerRetryTimeout(time.Millisecond),
		)
		checkResult(t, context.DeadlineExceeded, err, 3, attempt)
	})

	t.Run("TimeoutError2", func(t *testing.T) {
		var attempt int
		fn := func(ctx context.Context) error {
			attempt++

			timer := time.NewTimer(time.Millisecond * time.Duration(3-attempt))
			defer timer.Stop()
			select {
			case <-timer.C:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		err := DoWithContext(context.Background(), fn,
			WithMax(3),
			WithIsRetriable(IsDeadlineExceededError),
			WithPerRetryTimeout(time.Millisecond),
		)
		checkResult(t, nil, err, 2, attempt)
	})
}

func checkResult(t *testing.T, expectedErr, gotErr error, expectedAttempt, gotAttempt int) {
	if expectedErr != gotErr {
		t.Errorf("expected: %v, got: %v", expectedErr, gotErr)
	}
	if expectedAttempt != gotAttempt {
		t.Errorf("expected: %d, got: %d", expectedAttempt, gotAttempt)
	}
}

package retry

import (
	"context"
	"time"
)

type options struct {
	Max            uint
	PerCallTimeout time.Duration
	IsRetriable    func(error) bool
}

type RetryOption func(*options)

func WithMax(maxRetries uint) RetryOption {
	return func(o *options) {
		o.Max = maxRetries
	}
}

func WithPerRetryTimeout(timeout time.Duration) RetryOption {
	return func(o *options) {
		o.PerCallTimeout = timeout
	}
}

func WithIsRetriable(fn func(error) bool) RetryOption {
	return func(o *options) {
		o.IsRetriable = fn
	}
}

func Do(fn func() error, retryOptions ...RetryOption) error {
	// disable timeout control for function without ctx parameter
	retryOptions = append(retryOptions, WithPerRetryTimeout(0))
	return DoWithContext(context.Background(), func(context.Context) error {
		return fn()
	}, retryOptions...)
}

func DoWithContext(ctx context.Context, fn func(context.Context) error, retryOptions ...RetryOption) error {
	o := newOptions(retryOptions...)
	if o.Max == 0 {
		return fn(ctx)
	}
	var err error
	for i := uint(0); i < o.Max; i++ {
		callCtx, cancelFn := ctx, func() {}
		if o.PerCallTimeout > 0 {
			callCtx, cancelFn = context.WithTimeout(callCtx, o.PerCallTimeout)
		}
		defer cancelFn()
		err = fn(callCtx)
		if o.IsRetriable == nil || !o.IsRetriable(err) {
			return err
		}
	}

	return err

}

func newOptions(retryOptions ...RetryOption) *options {
	var o options
	for _, retryOption := range retryOptions {
		retryOption(&o)
	}
	return &o
}

func IsDeadlineExceededError(err error) bool {
	return err == context.DeadlineExceeded
}

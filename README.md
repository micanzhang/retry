# retry

golang version retry package, big influence by [github.com/grpc-ecosystem/go-grpc-middleware/grpc_retry](https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/retry/retry.go). 


## Usage

```go

import "github.com/micanzhang/retry"

func foo() error {
	return DoWithContext(context.Background(), fn,
		WithMax(3),
		WithIsRetriable(retry.IsDeadlineExceededError),
		WithPerRetryTimeout(time.Millisecond))
}
```

## Options

### WithMax

maximum number of retries on this call.

### WithIsRetriable

function decide whether retry or not. for example retry if err not nil:

```go
func HasError(err error) bool {
    return err != nil
}
```

### WithPerRetryTimeout

timeout per on this call.

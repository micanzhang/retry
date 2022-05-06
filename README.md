# retry

[![Build Status](https://github.com/micanzhang/retry/workflows/Go/badge.svg?branch=main)](https://github.com/micanzhang/retry/actions?query=branch=main)
[![codecov](https://codecov.io/gh/micanzhang/retry/branch/main/graph/badge.svg)](https://codecov.io/gh/micanzhang/retry)
[![GoDoc](https://pkg.go.dev/badge/github.com/micanzhang/retry?status.svg)](https://pkg.go.dev/github.com/micanzhang/retry?tab=doc)

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

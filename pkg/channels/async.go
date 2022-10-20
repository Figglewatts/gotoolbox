package channels

import (
	"context"
)

type AsyncResult[T any] struct {
	Value T
	Err   error
}

func (r AsyncResult[T]) IsError() bool {
	return r.Err != nil
}

// DoAsync can easily run a piece of non-preemptive code
// in a goroutine, allowing you to wait for the result
// and see if any errors occurred.
func DoAsync[T any](
	ctx context.Context, f func() (T, error),
) <-chan AsyncResult[T] {
	resultChan := make(chan AsyncResult[T])
	go func() {
		defer close(resultChan)
		result, err := f()
		select {
		case <-ctx.Done():
			return
		case resultChan <- AsyncResult[T]{Err: err, Value: result}:
		}
	}()
	return resultChan
}

package channels

import (
	"context"
)

// Receive will receive all values of the given input channel and append each
// to a slice, returning it.
// Receive will block until all values have been read from the given channel.
//  channel := Generate(ctx, 1, 2, 3, 4)
//  vals := Receive(ctx, channel)
func Receive[T any](ctx context.Context, channel <-chan T) []T {
	result := make([]T, 0)
	for val := range ContextGuard(ctx, channel) {
		result = append(result, val)
	}
	return result
}

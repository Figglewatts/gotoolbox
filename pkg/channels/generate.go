package channels

import (
	"context"
)

// Generate sends the variadic values provided to a channel of the same type,
// closing the channel after all values are sent.
//  channel := Generate(ctx, 1, 2, 3, 4)
func Generate[T any](ctx context.Context, vals ...T) <-chan T {
	result := make(chan T)
	go func() {
		defer close(result)
		for _, v := range vals {
			select {
			case <-ctx.Done():
				return
			case result <- v:
			}
		}
	}()
	return result
}

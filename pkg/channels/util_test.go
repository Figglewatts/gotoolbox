package channels

import (
	"context"
)

// cancelAfterOne reads one value from the given channel
// and calls the [context.CancelFunc] on the next iteration
func cancelAfterOne[T any](
	channel <-chan T, cancel context.CancelFunc,
) <-chan T {
	results := make(chan T)
	go func() {
		defer close(results)
		i := 0
		for val := range channel {
			if i == 1 {
				cancel()
				return
			}
			results <- val
			i++
		}
	}()
	return results
}

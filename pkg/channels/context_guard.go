package channels

import (
	"context"
)

// ContextGuard reads values from the given channel into the returned
// channel, until the given channel is closed or the context is done.
//
// Intended to be used with a for...range loop like so:
//  for val := range ContextGuard(ctx, inputChannel) {}
func ContextGuard[T any](ctx context.Context, c <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			// guard against context done when reading from
			// the channel
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return // account for channel closing
				}
				select {
				case <-ctx.Done(): // also guard when writing
					return
				case valStream <- v:
				}
			}
		}
	}()
	return valStream
}

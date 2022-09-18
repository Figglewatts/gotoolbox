package channels

import (
	"context"
	"sync"
)

// Multiplex will combine multiple channels of the same type into one. Values
// sent to any of the given channels will be multiplexed to the return channel.
// Note that order of values received may not be preserved.
//  c1 := Generate(ctx, 1, 2, 3)
//  c2 := Generate(ctx, 4, 5, 6)
//  multiplexed := Multiplex(ctx, c1, c2)
func Multiplex[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	wg := sync.WaitGroup{}
	multiplexedStream := make(chan T)

	multiplexChannel := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplexChannel(c)
	}

	// wait for all to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

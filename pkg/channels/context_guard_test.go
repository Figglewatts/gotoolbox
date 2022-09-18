package channels

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextGuard(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	channel := make(chan int)
	go func() {
		for v := range []int{1, 2, 3, 4} {
			channel <- v
		}
	}()
	guardedChan := ContextGuard(ctx, channel)

	i := 0
	for _ = range guardedChan {
		i++
		if i == 1 {
			cancel()
		}
	}

	assert.Equal(t, 1, i, "Incorrect number of iterations")
}

func TestContextGuardClose(t *testing.T) {
	ctx := context.Background()
	sendCtx, sendCancel := context.WithCancel(ctx)
	defer sendCancel()

	channel := make(chan int)
	go func() {
		for v := range []int{1, 2, 3, 4} {
			select {
			case <-sendCtx.Done():
				return
			case channel <- v:
			}
		}
	}()
	guardedChan := ContextGuard(ctx, channel)

	i := 0
	for _ = range guardedChan {
		i++
		if i == 1 {
			sendCancel()
			close(channel)
		}
	}

	assert.Equal(t, 2, i, "Incorrect number of iterations")
}

func ExampleContextGuard() {
	ctx := context.Background()
	channel := Generate(ctx, 1, 2, 3, 4)

	for val := range ContextGuard(ctx, channel) {
		fmt.Println(val)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}

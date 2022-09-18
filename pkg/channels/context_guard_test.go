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
		defer close(channel)
		for v := range []int{1, 2, 3, 4} {
			channel <- v
		}
	}()
	guardedChan := ContextGuard(ctx, channel)

	i := 0
	for _ = range guardedChan {
		cancel()
		i++
	}

	assert.Less(t, i, 4, "Iterations should be less than input length")
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

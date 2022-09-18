package channels

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplex(t *testing.T) {
	ctx := context.Background()

	a := Generate(ctx, 0, 2, 4, 6, 8)
	b := Generate(ctx, 1, 3, 5, 7, 9)
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	results := Multiplex(ctx, a, b)

	received := Receive(ctx, results)
	assert.ElementsMatch(t, expected, received)
}

func TestMultiplexCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := Generate(ctx, 0, 2)
	b := Generate(ctx, 1)

	results := Multiplex(ctx, a, b)

	i := 0
	for range results {
		i++
		if i == 2 {
			cancel()
		}
	}

	assert.Less(t, i, 3, "Iterations must be less than input size")
}

func ExampleMultiplex() {
	ctx := context.Background()
	c1 := Generate(ctx, 1, 2, 3)
	c2 := Generate(ctx, 4, 5, 6)

	multiplexed := Multiplex(ctx, c1, c2)

	received := Receive(ctx, multiplexed)
	sort.Ints(received)
	fmt.Println(received)
	// Output: [1 2 3 4 5 6]
}

package channels

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	ctx := context.Background()

	generated := Generate(ctx, 1, 2, 3, 4)

	expected := []int{1, 2, 3, 4}
	result := Receive(ctx, generated)

	assert.Equal(t, expected, result)
}

func TestGenerateCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generated := Generate(ctx, 1, 2, 3, 4)

	i := 0
	for range generated {
		i++
		if i == 1 {
			cancel()
		}
	}

	assert.Less(t, i, 4, "Iterations must be less than input size")
}

func ExampleGenerate() {
	ctx := context.Background()
	channel := Generate(ctx, 1, 2, 3, 4)

	received := Receive(ctx, channel)
	fmt.Println(received)

	// Output: [1 2 3 4]
}

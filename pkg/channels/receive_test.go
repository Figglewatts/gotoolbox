package channels

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceive(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	ctx := context.Background()
	inputChan := Generate(ctx, expected...)

	result := Receive(ctx, inputChan)

	assert.Equal(t, expected, result)
}

func ExampleReceive() {
	ctx := context.Background()
	channel := Generate(ctx, 1, 2, 3, 4)
	vals := Receive(ctx, channel)
	fmt.Println(vals)
	// Output: [1 2 3 4]
}

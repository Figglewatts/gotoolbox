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

	v, ok := <-generated
	if ok {
		t.Errorf("generated channel was not closed, got: %v", v)
	}
}

func TestGenerateCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	channel := cancelAfterOne(Generate(ctx, 1, 2, 3, 4), cancel)

	// read the first value
	<-channel

	v, ok := <-channel
	if ok {
		t.Errorf("channel should be closed, got: %v", v)
	}
}

func ExampleGenerate() {
	ctx := context.Background()
	channel := Generate(ctx, 1, 2, 3, 4)

	received := Receive(ctx, channel)
	fmt.Println(received)

	// Output: [1 2 3 4]
}

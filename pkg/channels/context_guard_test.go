package channels

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextGuard(t *testing.T) {
	ctx := context.Background()

	expected := []int{1, 2, 3, 4}
	channel := Generate(ctx, expected...)
	guardedChannel := ContextGuard(ctx, channel)

	result := Receive(ctx, guardedChannel)

	assert.Equal(t, expected, result)

	v, ok := <-guardedChannel
	if ok {
		t.Errorf("guardedChannel was not closed, got: %v", v)
	}
}

func TestContextGuardCancelBeforeRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	channel := make(chan int)
	guardedChannel := ContextGuard(ctx, channel)
	go func() {
		defer close(channel)

		// send a value on the input channel to unblock the read
		channel <- 1

		// now cancel the context right before we unblock the
		// channel write by reading a value in the test - hopefully triggering
		// the 2nd context check
		cancel()
	}()

	// wait a little bit before reading, to let the context settle
	<-time.After(time.Millisecond * 1)

	v, ok := <-guardedChannel
	if ok {
		t.Errorf("guardedChannel was not closed, got: %v", v)
	}
}

func TestContextGuardCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel, triggering the 1st context check

	channel := make(chan int)
	defer close(channel)
	guardedChannel := ContextGuard(ctx, channel)

	// reading from the unguarded channel should block as no values
	// will be sent
	select {
	case v := <-channel:
		t.Errorf("<-channel == %v want nothing (it should block)", v)
	default:
	}

	// however reading from the guarded channel will not block, as
	// we are guarding for the context being done
	// there should not be any values on the channel
	for val := range guardedChannel {
		t.Errorf("<-guardedChannel == %v want nothing (it should close)", val)
	}

	// and the channel should be closed
	v, ok := <-guardedChannel
	if ok {
		t.Errorf("<-guardedChannel should be closed but got %v", v)
	}
}

func TestContextGuardCancelPartial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	channel := cancelAfterOne(Generate(ctx, 1, 2, 3, 4), cancel)
	guarded := ContextGuard(ctx, channel)

	// read the first value
	<-guarded

	v, ok := <-guarded
	if ok {
		t.Errorf("guarded channel should be closed, got: %v", v)
	}
}

func ExampleContextGuard() {
	channel := Generate(context.Background(), 1, 2, 3)

	ctx := context.Background()
	for val := range ContextGuard(ctx, channel) {
		fmt.Println(val)
	}

	// Output:
	// 1
	// 2
	// 3
}

package channels

import (
	"context"
	"fmt"
	"testing"
)

func TestContextGuard(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel

	channel := make(chan int)
	guardedChannel := ContextGuard(ctx, channel)

	// reading from the unguarded channel should block as no values
	// will be sent
	select {
	case v := <-channel:
		t.Errorf("<-channel == %v want nothing (it should block)", v)
	default:
	}

	// however reading from the guarded channel will not block, as
	// we are guarding for the context being done - we should not
	// perform any iterations here
	for val := range guardedChannel {
		t.Errorf("<-guardedChannel == %v want nothing (it should close)", val)
	}
}

func TestContextGuardCancel(t *testing.T) {
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

func TestExampleContextGuard(t *testing.T) {
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

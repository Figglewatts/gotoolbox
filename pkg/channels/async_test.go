package channels

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleDoAsync() {
	ctx := context.Background()
	result := DoAsync(
		ctx, func() (string, error) {
			// simulate something taking some time to compute
			<-time.After(time.Millisecond * 3)
			return "some result", nil
		},
	)

	res := <-result
	fmt.Println(res.Value)
	// Output: some result
}

func TestDoAsync(t *testing.T) {
	ctx := context.Background()
	doSomething := func(a, b string) (string, error) {
		return a + b, nil
	}

	// do something that should return a value
	expected := "ab"
	resultChan := DoAsync(
		ctx, func() (string, error) {
			return doSomething("a", "b")
		},
	)
	res := <-resultChan
	if res.IsError() {
		t.Errorf("<-resultChan returned an error: %v", res)
	}
	assert.Equal(t, expected, res.Value)

	doSomethingErr := func(a, b string) (string, error) {
		return "", errors.New("something terrible has happened")
	}

	// now do something that should return an error
	resultChan = DoAsync(
		ctx, func() (string, error) {
			return doSomethingErr("a", "b")
		},
	)
	res = <-resultChan
	if !res.IsError() {
		t.Errorf(
			"<-resultChan should have returned an error, got: %v", res.Value,
		)
	}
}

func TestDoAsyncCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := DoAsync(
		ctx, func() (string, error) {
			cancel() // cancel before we've finished 'doing the work'
			return "result", nil
		},
	)

	// wait a little bit before reading, to let the context settle
	<-time.After(1 * time.Millisecond)

	for val := range resultChan {
		t.Errorf("<-resultChan should be closed, got: %v", val.Value)
	}
}

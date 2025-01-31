package operation_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/openmeterio/openmeter/pkg/framework/operation"
	"github.com/stretchr/testify/assert"
)

type ExampleRequest struct {
	Name string
}

type ExampleResponse struct {
	Greeting string
}

func exampleOperation(ctx context.Context, request ExampleRequest) (ExampleResponse, error) {
	if request.Name == "" {
		return ExampleResponse{}, errors.New("name is required")
	}
	return ExampleResponse{Greeting: "Hello, " + request.Name}, nil
}

func ExampleOperation() {
	var op operation.Operation[ExampleRequest, ExampleResponse] = exampleOperation

	resp, err := op(context.Background(), ExampleRequest{Name: "World"})
	if err != nil {
		panic(err)
	}

	fmt.Print(resp.Greeting)
	// Output: Hello, World
}

func TestAsNoResponseOperation(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		called := false
		f := func(ctx context.Context, req string) error {
			called = true
			return nil
		}

		op := operation.AsNoResponseOperation(f)
		resp, err := op(context.Background(), "test")

		assert.Nil(t, err)
		assert.Nil(t, resp)
		assert.True(t, called)
	})

	t.Run("error case", func(t *testing.T) {
		expectedErr := errors.New("test error")
		f := func(ctx context.Context, req string) error {
			return expectedErr
		}

		op := operation.AsNoResponseOperation(f)
		resp, err := op(context.Background(), "test")

		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)
	})
}

func TestCompose(t *testing.T) {
	t.Run("successful composition", func(t *testing.T) {
		op1 := func(ctx context.Context, req string) (int, error) {
			return len(req), nil
		}

		op2 := func(ctx context.Context, req int) (bool, error) {
			return req > 5, nil
		}

		composed := operation.Compose(op1, op2)
		result, err := composed(context.Background(), "hello world")

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("first operation fails", func(t *testing.T) {
		expectedErr := errors.New("op1 failed")
		op1 := func(ctx context.Context, req string) (int, error) {
			return 0, expectedErr
		}

		op2 := func(ctx context.Context, req int) (bool, error) {
			return true, nil
		}

		composed := operation.Compose(op1, op2)
		result, err := composed(context.Background(), "test")

		assert.Equal(t, expectedErr, err)
		assert.False(t, result)
	})

	t.Run("second operation fails", func(t *testing.T) {
		expectedErr := errors.New("op2 failed")
		op1 := func(ctx context.Context, req string) (int, error) {
			return 42, nil
		}

		op2 := func(ctx context.Context, req int) (bool, error) {
			return false, expectedErr
		}

		composed := operation.Compose(op1, op2)
		result, err := composed(context.Background(), "test")

		assert.Equal(t, expectedErr, err)
		assert.False(t, result)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		op1 := func(ctx context.Context, req string) (int, error) {
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			return 42, nil
		}

		op2 := func(ctx context.Context, req int) (bool, error) {
			return true, nil
		}

		composed := operation.Compose(op1, op2)
		_, err := composed(ctx, "test")

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

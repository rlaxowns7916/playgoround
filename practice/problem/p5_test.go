package problem

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestP5Drain(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(ctx, 5)
	counter := atomic.Int64{}

	for i := 0; i < 100; i++ {
		require.NoError(t, pool.Submit(ctx, func(context.Context) error { counter.Add(1); return nil }))
	}
	pool.Close()
	waitErr := pool.Wait()

	if waitErr != nil {
		t.Fatalf("wait error: %v", waitErr)
	}
	if counter.Load() != 100 {
		t.Fatalf("counter is %d, want 100", counter.Load())
	}
}

func TestP5WithRootContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := NewPool(ctx, 5)

	for i := 0; i < 100; i++ {
		require.NoError(t, pool.Submit(ctx, func(context.Context) error { return nil }))
	}
	cancel()
	waitErr := pool.Wait()

	if waitErr == nil {
		t.Fatalf("wait shoud return canceled error: %v", waitErr)
	}
}

func TestP5WithDenySubmitAfterClose(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(ctx, 5)

	for i := 0; i < 100; i++ {
		require.NoError(t, pool.Submit(ctx, func(context.Context) error { return nil }))
	}

	pool.Close()
	submitErr := pool.Submit(ctx, func(context.Context) error { return nil })

	if submitErr == nil {
		t.Fatalf("submit after close shoud return error: %v", submitErr)
	}
}

func TestP5SubmitCanCanceledWithSubmitContext(t *testing.T) {
	rootCtx := context.Background()
	pool := NewPool(rootCtx, 0)

	submitCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	submitErr := pool.Submit(submitCtx, func(context.Context) error {
		return nil
	})

	if submitErr != nil {
		t.Fatalf("submit shoud return nil error: %v", submitErr)
	}

	if errors.Is(submitErr, ErrSubmitCanceled) {
		t.Fatalf("submit shoud canceled error: %v", submitErr)
	}
}

func TestP5WithMultipleCloseIsSafe(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(ctx, 5)

	pool.Close()
	pool.Close()
	pool.Close()
}

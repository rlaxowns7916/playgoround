package problem

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithNonBuffer(t *testing.T) {
	ctx := context.Background()
	history, err := Rally(ctx, 2000, 0)

	if err != nil {
		t.Fatalf("Rally throws error: %v", err)
	}

	if len(history) != 2000*2 {
		t.Fatalf("history length is %v, want %d", len(history), 2000*2)
	}

	for i := 0; i < 2000*2-1; i++ {
		if history[i] == history[i+1] {
			t.Fatalf("history[%d] == history[%d]", i, i)
		}
	}
}

func TestWithBuffer(t *testing.T) {
	ctx := context.Background()
	history, err := Rally(ctx, 2000, 10)

	if err != nil {
		t.Fatalf("Rally throws error: %v", err)
	}

	if len(history) != 2000*2 {
		t.Fatalf("history length is %v, want %d", len(history), 2000*2)
	}

	for i := 0; i < 2000*2-1; i++ {
		if history[i] == history[i+1] {
			t.Fatalf("history[%d] == history[%d]", i, i)
		}
	}
}

func TestRallyWithTimeOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := Rally(ctx, 2000000, 10)

	if err == nil {
		t.Fatalf("Rally should not be nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Rally should have timed out")
	}
}

func TestRallyWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Rally(ctx, 2000000, 0)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected err : %s", err)
	}
}

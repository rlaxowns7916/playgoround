package problem

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

func TestP4WithComplete(t *testing.T) {
	ctx := context.Background()
	n := 5000
	buf := 8
	w := 4

	results, err := Pipeline(ctx, n, buf, w)

	if err != nil {
		t.Fatalf("expect no error")
	}

	for _, result := range results {
		if result%3 != 0 {
			t.Fatalf("non-multiple of 3 passed: %d", result)
		}

		root := math.Sqrt(float64(result))
		if root != math.Trunc(root) {
			t.Fatalf("non-square value passed: %d", result)
		}
	}
}

func TestP4WithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	n := 5000
	buf := 8
	w := 4

	cancel()
	_, err := Pipeline(ctx, n, buf, w)

	if err == nil {
		t.Fatalf("expect error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expect context.Canceled error")
	}
}

func TestP4WithTimeOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	n := 1_000_000_000
	buf := 0
	w := 1

	_, err := Pipeline(ctx, n, buf, w)

	if err == nil {
		t.Fatalf("expect error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expect context.DeadLineExceed error")
	}
}

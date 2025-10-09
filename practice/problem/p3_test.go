package problem

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestP3With1Worker(t *testing.T) {
	n := 30000
	w := 1
	buf := 100
	work := func(task int) int {
		return task * task
	}

	ctx := context.Background()
	results, err := Run(ctx, n, w, buf, work)

	if err != nil {
		t.Fatalf("expect an no error")
	}

	if n != len(results) {
		t.Fatalf("expect results length %d, but got %d", n, len(results))
	}
}

func TestP3WithNWorker(t *testing.T) {
	n := 30000
	w := 100
	buf := 100
	work := func(task int) int {
		return task * task
	}

	ctx := context.Background()
	results, err := Run(ctx, n, w, buf, work)

	if err != nil {
		t.Fatalf("expect no error")
	}

	if n != len(results) {
		t.Fatalf("expect results length %d, but got %d", n, len(results))
	}
}

func TestP3WithBuffer(t *testing.T) {
	n := 30000
	w := 1
	buf := 1000
	work := func(task int) int {
		return task * task
	}

	ctx := context.Background()
	results, err := Run(ctx, n, w, buf, work)

	if err != nil {
		t.Fatalf("expect an no error")
	}

	if n != len(results) {
		t.Fatalf("expect results length %d, but got %d", n, len(results))
	}
}

func TestP3WithNonBuffer(t *testing.T) {
	n := 30000
	w := 1
	buf := 0
	work := func(task int) int {
		return task * task
	}

	ctx := context.Background()
	results, err := Run(ctx, n, w, buf, work)

	if err != nil {
		t.Fatalf("expect an no error")
	}

	if n != len(results) {
		t.Fatalf("expect results length %d, but got %d", n, len(results))
	}
}

func TestP3WithTimeOut(t *testing.T) {
	n := 30000
	w := 10
	buf := 100
	work := func(task int) int {
		return task * task
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := Run(ctx, n, w, buf, work)

	if err == nil {
		t.Fatalf("expect an error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expect DeadlineExceeded error")
	}
}

func TestP3WithCancel(t *testing.T) {
	n := 30000
	w := 10
	buf := 100
	work := func(task int) int {
		return task * task
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Run(ctx, n, w, buf, work)

	if err == nil {
		t.Fatalf("expect an error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expect Canceled error")
	}
}

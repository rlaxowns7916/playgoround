package problem

import (
	"context"
	"testing"
)

func TestP2WithNonBuffer(t *testing.T) {
	n := 30
	ctx := context.Background()
	inputs := make([]<-chan int, 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan int, 3)
		ch <- i
		ch <- i + 1
		ch <- i + 2

		inputs = append(inputs, ch)
		close(ch)
	}

	answer, err := FanIn(ctx, inputs, 0)

	if err != nil {
		t.Fatalf("FanIn returned an error: %v", err)
	}

	if len(answer) != n*3 {
		t.Fatalf("FanIn returned %d answers, want %d", len(answer), n*3)
	}
}

func TestP2WithBuffer(t *testing.T) {
	n := 30
	ctx := context.Background()
	inputs := make([]<-chan int, 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan int, 3)
		ch <- i
		ch <- i + 1
		ch <- i + 2

		inputs = append(inputs, ch)
		close(ch)
	}

	answer, err := FanIn(ctx, inputs, 10)

	if err != nil {
		t.Fatalf("FanIn returned an error: %v", err)
	}

	if len(answer) != n*3 {
		t.Fatalf("FanIn returned %d answers, want %d", len(answer), n*3)
	}
}

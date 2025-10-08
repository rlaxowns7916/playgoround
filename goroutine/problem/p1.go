package problem

import (
	"context"
	"fmt"
	"sync"
)

type token = struct{}

func Rally(ctx context.Context, n, bufSize int) ([]string, error) {
	var wg sync.WaitGroup
	pingCh := make(chan token, bufSize)
	pongCh := make(chan token, bufSize)
	events := make(chan string, 2*n)
	history := make([]string, 0, 2*n)

	wg.Add(2)
	go runPing(ctx, n, pingCh, pongCh, events, &wg)
	go runPong(ctx, pingCh, pongCh, events, &wg)
	go func() { wg.Wait(); close(events) }()

	for e := range events {
		history = append(history, e)
	}
	if err := ctx.Err(); err != nil {
		return history, err
	}
	return history, nil
}

func runPing(ctx context.Context, n int, ping chan<- token, pong <-chan token, events chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ping)
	for i := 0; i < n; i++ {
		if err := send(ctx, ping, token{}); err != nil {
			fmt.Printf("runPing -> send signal: %v\n", err)
			return
		}

		if _, ok, err := receive(ctx, pong); err != nil || !ok {
			fmt.Printf("runPing -> receive signal: %v\n", err)
			return
		}

		if err := send(ctx, events, "pong"); err != nil {
			fmt.Printf("runPing -> send events: %v\n", err)
			return
		}
	}
}

func runPong(ctx context.Context, ping <-chan token, pong chan<- token, events chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if _, ok, err := receive(ctx, ping); err != nil || !ok {
			fmt.Printf("runPong -> receive signal: %v, isChannelOpen:%v \n", err, ok)
			return
		}

		if err := send(ctx, events, "ping"); err != nil {
			fmt.Printf("runPong -> send events: %v\n", err)
			return
		}

		if err := send(ctx, pong, token{}); err != nil {
			fmt.Printf("runPong -> send signal: %v\n", err)
			return
		}
	}
}

func send[T any](ctx context.Context, ch chan<- T, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- v:
		return nil
	}
}

func receive[T any](ctx context.Context, ch <-chan T) (T, bool, error) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, false, ctx.Err()
	case v, ok := <-ch:
		return v, ok, nil
	}
}

package problem

import (
	"context"
	"sync"
)

func FanIn(ctx context.Context, inputs []<-chan int, outBuf int) ([]int, error) {
	wg := sync.WaitGroup{}
	answer := make([]int, 0, len(inputs))
	outputs := make(chan int, outBuf)

	for _, input := range inputs {
		if input == nil {
			continue
		}
		wg.Add(1)
		go aggregate(ctx, &wg, input, outputs)
	}

	go func() {
		wg.Wait()
		close(outputs)
	}()

	for {
		select {
		case output, ok := <-outputs:
			if !ok {
				return answer, ctx.Err()
			}
			answer = append(answer, output)
		}
	}
}

func aggregate(ctx context.Context, wg *sync.WaitGroup, input <-chan int, output chan<- int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-input:
			if !ok {
				return
			}
			output <- message
		}
	}
}

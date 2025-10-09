package problem

import (
	"context"
	"sync"
)

func Run(ctx context.Context, n, w, buf int, f func(int) int) ([]int, error) {
	wg := sync.WaitGroup{}
	processed := make(chan int, buf)
	answer := make([]int, 0, n)

	tasks := produce(ctx, &wg, n, buf)
	for i := 0; i < w; i++ {
		wg.Add(1)
		go consume(ctx, &wg, tasks, processed, f)
	}

	go func() {
		wg.Wait()
		close(processed)
	}()

	for result := range processed {
		answer = append(answer, result)
	}

	return answer, ctx.Err()
}

func produce(ctx context.Context, wg *sync.WaitGroup, n, bufSize int) <-chan int {
	ch := make(chan int, bufSize)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)

		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()

	return ch
}

func consume(ctx context.Context, wg *sync.WaitGroup, input <-chan int, output chan<- int, f func(int) int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-input:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return

			// output queuer가 buffer를 모두 사용하고 Blocking 될떄를 방지
			case output <- f(task):
			}
		}
	}
}

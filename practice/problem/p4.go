package problem

import (
	"context"
	"sync"
)

func Pipeline(ctx context.Context, n, buf, w int) ([]int, error) {
	pipeLine0 := stage0(ctx, n, buf)
	pipeLine1 := stage1(ctx, w, buf, pipeLine0)
	pipeLine2 := stage2(ctx, buf, pipeLine1)
	var answer []int

	for filtered := range pipeLine2 {
		answer = append(answer, filtered)
	}

	return answer, ctx.Err()
}

func stage0(ctx context.Context, n, buf int) chan int {
	ch := make(chan int, buf)
	go func() {
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

func stage1(ctx context.Context, w, buf int, prevPipeLine <-chan int) chan int {
	ch := make(chan int, buf)
	wg := sync.WaitGroup{}

	go func() {
		defer close(ch)

		for i := 0; i < w; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case value, ok := <-prevPipeLine:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case ch <- value * value:
						}
					}
				}
			}()
		}

		wg.Wait()
	}()

	return ch

}

func stage2(ctx context.Context, buf int, prevChannel <-chan int) chan int {
	ch := make(chan int, buf)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-prevChannel:
				if !ok {
					return
				}
				if value%3 == 0 {
					select {
					case <-ctx.Done():
						return
					case ch <- value:
					}
				}
			}
		}
	}()

	return ch
}

package pipeline

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

func Execute(wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	extractCh := extract(ctx, wg, cancel, "./input.txt")
	transformCh := transform(ctx, wg, extractCh)
	load(ctx, wg, cancel, transformCh, "./output.txt")

	wg.Wait()
}

func extract(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, path string) <-chan string {
	out := make(chan string)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(out)

		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("open file %s err: %v\n", path, err)
			cancel()
			return
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if !sc.Scan() {
				if err := sc.Err(); err != nil && !errors.Is(err, io.EOF) {
					fmt.Println("scan err:", err)
					cancel()
				}
				return
			}

			out <- sc.Text()
		}
	}()

	return out
}

func transform(ctx context.Context, wg *sync.WaitGroup, in <-chan string) <-chan string {
	out := make(chan string)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- "------------------" + v + "------------------"
			}
		}
	}()

	return out
}

func load(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, in <-chan string, path string) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			fmt.Println("open file err:", err)
			cancel()
			return
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		defer w.Flush()

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				if _, err := w.WriteString(v + "\n"); err != nil {
					fmt.Println("write err:", err)
					cancel()
					return
				}
			}
		}
	}()
}

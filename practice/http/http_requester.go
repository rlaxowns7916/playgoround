package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

var requests = []string{
	"https://youtube.com",
	"https://naver.com",
	"https://google.com",
	"https://github.com",
}

func Execute(wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, request := range requests {
		wg.Add(1)
		go func(ctx context.Context, url string) {
			defer wg.Done()
			select {
			case resp := <-get(ctx, url):
				fmt.Printf("Got response from %s: %s\n", url, resp)
				cancel()
			case <-ctx.Done():
				fmt.Printf("Shut down received %s \n", url)
				return
			}
		}(ctx, request)
	}

	wg.Wait()
}

func get(ctx context.Context, url string) chan interface{} {
	ch := make(chan interface{}, 1)

	go func() {
		defer close(ch)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			ch <- err
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ch <- err
			return
		}
		ch <- resp
	}()

	return ch
}

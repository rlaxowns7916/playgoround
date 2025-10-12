package problem

import "sync"

func FanInToChannel(done <-chan struct{}, channels ...chan interface{}) chan interface{} {
	if len(channels) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	out := make(chan interface{})

	wg.Add(len(channels))
	for _, channel := range channels {
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case input, ok := <-channel:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case out <- input:
					}
				}
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		defer close(out)
	}()

	return out
}

func FanOutToChannel(done <-chan struct{}, inputs ...interface{}) []chan interface{} {
	out := make([]chan interface{}, 0, len(inputs))

	if len(inputs) == 0 {
		return out
	}

	for _, input := range inputs {
		fanOutChannel := make(chan interface{})

		go func(v interface{}, channel chan<- interface{}) {
			defer close(channel)
			select {
			case channel <- v:
			case <-done:
				return
			}
		}(input, fanOutChannel)

		out = append(out, fanOutChannel)
	}

	return out
}

func repeat(done <-chan struct{}, generator func() interface{}, n int) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case out <- generator():
			}
		}
	}()

	return out
}

func take(done <-chan struct{}, input <-chan interface{}, n int) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case value, ok := <-input:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}

func bridge(done <-chan struct{}, channelStream <-chan <-chan interface{}) <-chan interface{} {
	wg := sync.WaitGroup{}
	out := make(chan interface{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case channel, ok := <-channelStream:
				if !ok {
					return
				}
				wg.Add(1)
				go func(ch <-chan interface{}) {
					defer wg.Done()
					for {
						select {
						case <-done:
							return
						case value, ok := <-ch:
							if !ok {
								return
							}
							out <- value
						}
					}
				}(channel)
			}
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

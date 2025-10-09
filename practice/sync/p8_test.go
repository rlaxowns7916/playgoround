package sync

import (
	"sync/atomic"
	"testing"
)

func TestP8InitExactlyOnce(t *testing.T) {
	counter := atomic.Int64{}
	initFunc := func() (*Resource, error) {
		counter.Add(1)
		return &Resource{}, nil
	}

	lazy := NewLazy(initFunc)
	for i := 0; i < 100; i++ {
		go func() {
			_, _ = lazy.Get()
		}()
	}

	if counter.Load() != 1 {
		t.Fatalf("counter.Load() should be 1")
	}
}

package problem

import (
	"testing"
	"time"
)

func TestP9WhenChannelParameterEmpty(t *testing.T) {
	channel := Or()
	_, ok := <-channel

	if ok {
		t.Fatalf("channel is opened")
	}
}

func TestP9With2ParametersAndOneChannelIsClosed(t *testing.T) {
	chan1 := make(chan interface{})
	close(chan1)
	chan2 := make(chan interface{})

	channel := Or(chan1, chan2)
	select {
	case _, ok := <-channel:
		if ok {
			t.Fatalf("expected closed channel")
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("should close immediately")
	}
}

func TestP9WithNParametersAndOneChannelIsClosed(t *testing.T) {
	chan1 := make(chan interface{})
	close(chan1)
	chan2 := make(chan interface{})
	chan3 := make(chan interface{})
	chan4 := make(chan interface{})
	chan5 := make(chan interface{})

	channel := Or(chan1, chan2, chan3, chan4, chan5)

	select {
	case _, ok := <-channel:
		if ok {
			t.Fatalf("expected closed channel")
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("should close immediately")
	}
}

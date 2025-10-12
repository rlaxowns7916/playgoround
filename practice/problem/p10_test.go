package problem

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestP10Repeat(t *testing.T) {

	done := make(chan struct{})
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	randGenerator := func() interface{} {
		return random.Int()
	}

	result := make([]interface{}, 0, 0)
	for output := range repeat(done, randGenerator, 10) {
		result = append(result, output)
	}

	if len(result) != 10 {
		t.Fatalf("expect 10 but %d", len(result))
	}

	fmt.Println(result)
}

func TestP10Take(t *testing.T) {
	done := make(chan struct{})
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	randGenerator := func() interface{} {
		return random.Int()
	}

	result := make([]interface{}, 0, 0)
	for output := range take(done, repeat(done, randGenerator, 10), 4) {
		result = append(result, output)
	}
	close(done)

	if len(result) != 4 {
		t.Fatalf("expect 4 but %d", len(result))
	}

	fmt.Println(result)
}

func TestP10Bridge(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	channelStream := make(chan (<-chan interface{}))

	go func() {
		defer close(channelStream)
		for i := 0; i < 5; i++ {
			ch := repeat(done, func() interface{} { return i }, 1)
			channelStream <- ch
		}
	}()

	result := make([]interface{}, 0, 0)
	for out := range bridge(done, channelStream) {
		result = append(result, out)
	}

	if len(result) != 5 {
		t.Fatalf("expect 5 but %d", len(result))
	}

	fmt.Println(result)

}

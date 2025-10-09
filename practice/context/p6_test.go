package context

import (
	"context"
	"testing"
)

func TestP6ExistValue(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequest(ctx, "reqID", 1)
	requestID, ok := RequestID(ctx)
	if !ok {
		t.Fatalf("expect value")
	}
	if requestID != "reqID" {
		t.Fatalf("expect reqID")
	}
}

func TestP6NotExistValue(t *testing.T) {
	ctx := context.Background()
	requestID, ok := RequestID(ctx)

	if ok {
		t.Fatalf("expect non-value")
	}

	if requestID != "" {
		t.Fatalf("expect empty string")
	}
}

func TestP6WithChildContextCanGetParentContextValue(t *testing.T) {
	rootCtx := WithRequest(context.Background(), "rid", 42)
	childCtx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	requestId, ok := RequestID(childCtx)

	if !ok {
		t.Fatalf("expect value")
	}
	if requestId != "rid" {
		t.Fatalf("expect rid")
	}
}

func TestP6WithParentContextCannotGetChildContextValue(t *testing.T) {
	rootCtx := WithRequest(context.Background(), "rid", 42)
	context.WithValue(rootCtx, keyUserID, 1)

	userID, _ := UserID(rootCtx)

	if userID != 42 {
		t.Fatalf("expect userId 42 bot get : %d", userID)
	}
}

func TestP6WithPropagationInGoroutine(t *testing.T) {
	t.Parallel()
	ctx := WithRequest(context.Background(), "rid", 42)
	ch := make(chan string, 1)
	go func(c context.Context) {
		id, _ := RequestID(c)
		ch <- id
	}(ctx)
	got := <-ch
	if got != "rid" {
		t.Fatalf("want rid, got %s", got)
	}
}

package context

import "context"

type key int

const (
	keyReqID key = iota
	keyUserID
)

func WithRequest(parent context.Context, reqID string, userID int64) context.Context {
	ctx := context.WithValue(parent, keyReqID, reqID)
	return context.WithValue(ctx, keyUserID, userID)
}

func RequestID(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(keyReqID).(string)
	return reqID, ok
}
func UserID(ctx context.Context) (int64, bool) {
	UserID, ok := ctx.Value(keyUserID).(int64)
	return UserID, ok
}

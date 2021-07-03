package util

import "context"

const userIDContextKey = "UserIDContextKey"

func SetCurrentUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetCurrentUserID(ctx context.Context) uint {
	return ctx.Value(userIDContextKey).(uint)
}

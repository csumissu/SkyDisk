package util

import (
	"context"
)

const userIDContextKey = "UserIDContextKey"

func SetCurrentUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetCurrentUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uint)
	return userID, ok
}

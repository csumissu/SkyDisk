package redis

import (
	"context"
	"github.com/csumissu/SkyDisk/infra"
	"github.com/csumissu/SkyDisk/util/logger"
	"time"
)

var ctx = context.Background()
var redisClient = infra.RedisClient

func Set(key string, value interface{}, expiration time.Duration) bool {
	err := redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Error("set redis value failed, key: %s, value: %v", key, value, err)
		return false
	}
	return true
}

func Del(keys ...string) bool {
	err := redisClient.Del(ctx, keys...).Err()
	if err != nil {
		logger.Error("delete redis key failed, keys: %v", keys, err)
		return false
	}
	return true
}

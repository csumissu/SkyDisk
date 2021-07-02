package infra

import (
	"context"
	"fmt"
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/util"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisClient struct {
	client *redis.Client
}

var RedisClient *redisClient
var ctx = context.Background()

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisCfg.Host, config.RedisCfg.Port),
		Password: config.RedisCfg.Password,
		DB:       config.RedisCfg.DB,
	})

	RedisClient = &redisClient{client}
}

func (redis *redisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := redis.client.Set(ctx, key, value, expiration).Err(); err != nil {
		util.Logger.Error("set redis value failed, key: %v, value: %v", key, value, err)
		return false
	}
	return true
}

func (redis *redisClient) Del(keys ...string) bool {
	if err := redis.client.Del(ctx, keys...).Err(); err != nil {
		util.Logger.Error("delete redis key failed, keys: %v", keys, err)
		return false
	}
	return true
}

func (redis *redisClient) Exists(keys ...string) bool {
	if count, err := redis.client.Exists(ctx, keys...).Result(); err != nil {
		util.Logger.Error("check redis key exists failed, keys: %v", keys, err)
		return false
	} else {
		return count == int64(len(keys))
	}
}

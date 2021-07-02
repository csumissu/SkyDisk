package infra

import (
	"context"
	"fmt"
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/util"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

var instance *RedisClient
var ctx = context.Background()

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisCfg.Host, config.RedisCfg.Port),
		Password: config.RedisCfg.Password,
		DB:       config.RedisCfg.DB,
	})

	instance = &RedisClient{redisClient}
}

func Redis() *RedisClient {
	return instance
}

func (redis *RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	err := redis.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		util.Log().Error("set redis value failed, key: %s, value: %v", key, value, err)
		return false
	}
	return true
}

func (redis *RedisClient) Del(keys ...string) bool {
	err := redis.client.Del(ctx, keys...).Err()
	if err != nil {
		util.Log().Error("delete redis key failed, keys: %v", keys, err)
		return false
	}
	return true
}

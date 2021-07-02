package infra

import (
	"fmt"
	"github.com/csumissu/SkyDisk/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisCfg.Host, config.RedisCfg.Port),
		Password: config.RedisCfg.Password,
		DB:       config.RedisCfg.DB,
	})

	RedisClient = redisClient
}

package infrastructure

import (
	"fmt"
	"github.com/csumissu/SkyDisk/conf"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.RedisCfg.Host, conf.RedisCfg.Port),
		Password: conf.RedisCfg.Password,
		DB:       conf.RedisCfg.DB,
	})

	RedisClient = redisClient
}

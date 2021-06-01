package middleware

import (
	"fmt"
	"github.com/csumissu/SkyDisk/conf"
	"github.com/csumissu/SkyDisk/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Session() gin.HandlerFunc {
	address := conf.RedisCfg.Host + ":" + strconv.Itoa(conf.RedisCfg.Port)
	store, err := redis.NewStoreWithDB(10, conf.RedisCfg.Network, address, conf.RedisCfg.Password, conf.RedisCfg.DB)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to redis! %v", err))
	}

	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 24 * 3600, Path: "/"})
	return sessions.Sessions("sky-disk-sessions", store)
}

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid != nil {
			user, err := model.GetActiveUserByID(uid)
			if err == nil {
				c.Set("user", user)
			}
		}
		c.Next()
	}
}

func SetCurrentUser(c *gin.Context, user model.User) {
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	err := session.Save()
	if err != nil {
		log.Println("Set current user failed", err)
	}
}
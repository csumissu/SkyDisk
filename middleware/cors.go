package middleware

import (
	"github.com/csumissu/SkyDisk/conf"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     conf.CORSCfg.AllowOrigins,
		AllowMethods:     conf.CORSCfg.AllowMethods,
		AllowHeaders:     conf.CORSCfg.AllowHeaders,
		AllowCredentials: conf.CORSCfg.AllowCredentials,
		ExposeHeaders:    conf.CORSCfg.ExposeHeaders,
	})
}

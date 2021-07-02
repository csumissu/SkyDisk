package middleware

import (
	"github.com/csumissu/SkyDisk/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.CORSCfg.AllowOrigins,
		AllowMethods:     config.CORSCfg.AllowMethods,
		AllowHeaders:     config.CORSCfg.AllowHeaders,
		AllowCredentials: config.CORSCfg.AllowCredentials,
		ExposeHeaders:    config.CORSCfg.ExposeHeaders,
	})
}

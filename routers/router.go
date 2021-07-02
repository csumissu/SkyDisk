package routers

import (
	"github.com/csumissu/SkyDisk/middleware"
	"github.com/csumissu/SkyDisk/routers/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.GinContextToContextMiddleware())

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)
		api.POST("/graphql", middleware.AuthRequired(), controllers.GraphqlHandler())
	}

	return r
}

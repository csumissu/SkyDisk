package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/graph/generated"
	"github.com/csumissu/SkyDisk/middleware"
	"github.com/csumissu/SkyDisk/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var resolver = &Resolver{}

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.GinContextToContextMiddleware())

	api := r.Group("/api")
	{
		api.POST("/login", userLogin)
		api.POST("/graphql", graphqlHandler())
	}

	return r
}

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func userLogin(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err == nil {
		response := resolver.loginService.Login(c, request)
		c.JSON(response.HttpStatus, response)
	} else {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(err))
	}
}

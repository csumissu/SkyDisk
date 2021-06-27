package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/graph/generated"
	"github.com/csumissu/SkyDisk/middleware"
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
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := resolver.loginService.Login(c, request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

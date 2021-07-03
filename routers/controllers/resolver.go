package controllers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/csumissu/SkyDisk/routers/graph"
	"github.com/csumissu/SkyDisk/service"
	"github.com/gin-gonic/gin"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service.AuthService
	service.UserService
}

var r = Resolver{}

func GraphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &r}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

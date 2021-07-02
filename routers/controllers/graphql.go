package controllers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/csumissu/SkyDisk/routers/graphql/generated"
	"github.com/csumissu/SkyDisk/routers/graphql/resolvers"
	"github.com/gin-gonic/gin"
)

func GraphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

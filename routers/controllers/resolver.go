package controllers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/csumissu/SkyDisk/config"
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
	service.FileService
}

var r = Resolver{}

const MB = 1 << 20

func GraphqlHandler() gin.HandlerFunc {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &r}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     config.ServerCfg.MultipartMaxMemoryMB * MB,
		MaxUploadSize: config.ServerCfg.MultipartMaxUploadMB * MB,
	})
	srv.Use(extension.Introspection{})
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

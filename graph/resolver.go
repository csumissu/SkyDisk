package graph

import "github.com/csumissu/SkyDisk/service/auth"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	loginService auth.LoginService
}

package graph

import "github.com/csumissu/SkyDisk/service/user"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	authService user.AuthService
}

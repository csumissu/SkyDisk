package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/graph/generated"
)

func (r *mutationResolver) Login(ctx context.Context, input dto.LoginRequest) (*dto.LoginResponse, error) {
	return r.loginService.Login(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

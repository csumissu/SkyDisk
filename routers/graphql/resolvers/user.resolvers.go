package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/routers/graphql/generated"
)

func (r *queryResolver) SearchUserProfile(ctx context.Context) (*dto.UserProfileResponse, error) {
	return r.UserService.SearchUserProfile(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/csumissu/SkyDisk/graph/dto"
	"github.com/csumissu/SkyDisk/graph/generated"
)

func (r *queryResolver) SearchUserProfile(ctx context.Context) (*dto.UserProfileResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

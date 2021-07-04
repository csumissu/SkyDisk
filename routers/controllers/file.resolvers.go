package controllers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/routers/graph"
)

func (r *mutationResolver) SingleUpload(ctx context.Context, path string, file graphql.Upload) (bool, error) {
	return r.FileService.SingleUpload(ctx, path, file)
}

func (r *queryResolver) ListObjects(ctx context.Context, path string) ([]*dto.ObjectResponse, error) {
	return r.FileService.ListObjects(ctx, path)
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

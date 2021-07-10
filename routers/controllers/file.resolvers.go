package controllers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/routers/dto"
)

func (r *mutationResolver) SingleUpload(ctx context.Context, path string, file graphql.Upload) (bool, error) {
	return r.FileService.SingleUpload(ctx, path, file)
}

func (r *mutationResolver) DeleteObject(ctx context.Context, objectID uint) (bool, error) {
	return r.FileService.DeleteObject(ctx, objectID)
}

func (r *queryResolver) ListObjects(ctx context.Context, path string) (*dto.ListObjectsRresponse, error) {
	return r.FileService.ListObjects(ctx, path)
}

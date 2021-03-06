package controllers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/routers/dto"
)

func (r *mutationResolver) UploadFile(ctx context.Context, path string, file graphql.Upload) (bool, error) {
	return r.FileService.UploadFile(ctx, path, file)
}

func (r *mutationResolver) DeleteObject(ctx context.Context, objectID uint) (bool, error) {
	return r.FileService.DeleteObject(ctx, objectID)
}

func (r *mutationResolver) CreateDir(ctx context.Context, path string) (bool, error) {
	return r.FileService.CreateDir(ctx, path)
}

func (r *mutationResolver) RenameObject(ctx context.Context, objectID uint, newName string) (bool, error) {
	return r.FileService.RenameObject(ctx, objectID, newName)
}

func (r *mutationResolver) MoveObject(ctx context.Context, objectID uint, path string) (bool, error) {
	return r.FileService.MoveObject(ctx, objectID, path)
}

func (r *queryResolver) ListObjects(ctx context.Context, path string) (*dto.ListObjectsRresponse, error) {
	return r.FileService.ListObjects(ctx, path)
}

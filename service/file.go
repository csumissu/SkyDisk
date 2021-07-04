package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/routers/dto"
)

type FileService struct {
}

func (service *FileService) SingleUpload(ctx context.Context, path string, file graphql.Upload) (bool, error) {
	return true, nil
}

func (service *FileService) ListObjects(ctx context.Context, path string) ([]*dto.ObjectResponse, error) {
	return make([]*dto.ObjectResponse, 0), nil
}

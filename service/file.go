package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/infra/filesystem"
	"github.com/csumissu/SkyDisk/routers/dto"
)

type FileService struct {
}

func (service *FileService) SingleUpload(ctx context.Context, path string, upload graphql.Upload) (bool, error) {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return false, err
	}

	fs, err := filesystem.NewFileSystem(user)
	if err != nil {
		return false, err
	}
	defer fs.Recycle()

	fileInfo := filesystem.FileInfo{
		File:        upload.File,
		Name:        upload.Filename,
		Size:        uint64(upload.Size),
		MIMEType:    upload.ContentType,
		VirtualPath: path,
	}

	fs.Use(filesystem.HookAfterUpload, GenericAfterUpload)
	err = fs.Upload(ctx, fileInfo)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service *FileService) ListObjects(ctx context.Context, path string) ([]*dto.ObjectResponse, error) {
	return make([]*dto.ObjectResponse, 0), nil
}

func GenericAfterUpload(ctx context.Context, fs *filesystem.FileSystem) error {
	return nil
}

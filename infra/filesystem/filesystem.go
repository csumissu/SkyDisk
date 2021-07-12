package filesystem

import (
	"context"
	"fmt"
	"github.com/csumissu/SkyDisk/infra/filesystem/handlers/local"
	"github.com/csumissu/SkyDisk/models"
	"io"
	"sync"
)

type FileSystem struct {
	User    *models.User
	handler Handler
	hooks   map[string][]Hook
	mutex   sync.RWMutex
}

func NewFileSystem(user *models.User) (*FileSystem, error) {
	fs := getEmptyFS()
	fs.User = user
	fs.determineHandler()
	return fs, nil
}

func (fs *FileSystem) determineHandler() {
	fs.handler = local.Handler{}
}

func (fs *FileSystem) Upload(ctx context.Context, file io.Reader, info UploadFileInfo) error {
	ctx = context.WithValue(ctx, UploadFileInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath)

	if err := fs.handler.Put(ctx, file, objectKey, info.Size); err != nil {
		return err
	}

	return fs.Trigger(ctx, HookAfterUploadFile)
}

func (fs *FileSystem) Download(ctx context.Context, info DownloadObjectInfo) (io.ReadSeekCloser, error) {
	ctx = context.WithValue(ctx, DownloadObjectInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath)
	return fs.handler.Get(ctx, objectKey, info.IsDir)
}

func (fs *FileSystem) Delete(ctx context.Context, info DeleteObjectInfo) error {
	ctx = context.WithValue(ctx, DeleteObjectInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath)

	if err := fs.handler.Delete(ctx, objectKey, info.IsDir); err != nil {
		return err
	}

	return fs.Trigger(ctx, HookAfterDeleteObject)
}

func (fs *FileSystem) CreateDir(ctx context.Context, virtualPath string) error {
	ctx = context.WithValue(ctx, CreateDirCtx, virtualPath)
	objectKey := fs.generateObjectKey(virtualPath)

	if err := fs.handler.CreateDir(ctx, objectKey); err != nil {
		return err
	}

	return fs.Trigger(ctx, HookAfterCreateDir)
}

func (fs *FileSystem) generateObjectKey(virtualPath string) string {
	return fmt.Sprintf("uploads/%d/%s", fs.User.ID, virtualPath)
}

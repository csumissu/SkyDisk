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

func (fs *FileSystem) Upload(ctx context.Context, file io.Reader, info UploadObjectInfo) error {
	ctx = context.WithValue(ctx, UploadObjectInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath)

	if err := fs.handler.Put(ctx, file, objectKey, info.Size); err != nil {
		return err
	}

	return fs.Trigger(ctx, HookAfterUpload)
}

func (fs *FileSystem) Download(ctx context.Context, info DownloadObjectInfo) (io.ReadSeekCloser, error) {
	ctx = context.WithValue(ctx, DownloadObjectInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath)
	return fs.handler.Get(ctx, objectKey)
}

func (fs *FileSystem) Delete(ctx context.Context, info DeleteObjectInfo) error {
	ctx = context.WithValue(ctx, DeleteObjectInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath)

	if err := fs.handler.Delete(ctx, objectKey); err != nil {
		return err
	}

	return fs.Trigger(ctx, HookAfterDelete)
}

func (fs *FileSystem) generateObjectKey(virtualPath string) string {
	return fmt.Sprintf("uploads/%d/%s", fs.User.ID, virtualPath)
}

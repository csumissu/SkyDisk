package filesystem

import (
	"context"
	"fmt"
	"github.com/csumissu/SkyDisk/infra/filesystem/handlers/local"
	"github.com/csumissu/SkyDisk/models"
	"io"
	"path"
	"sync"
)

type FileSystem struct {
	User    *models.User
	handler Handler
	hooks   map[string][]Hook
	mutex   sync.RWMutex
}

type UploadFileInfo struct {
	File        io.Reader
	Name        string
	Size        uint64
	MIMEType    string
	VirtualPath string
}

type DownloadFileInfo struct {
	Name        *string
	VirtualPath string
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

func (fs *FileSystem) Upload(ctx context.Context, info UploadFileInfo) error {
	ctx = context.WithValue(ctx, UploadFileInfoCtx, info)
	objectKey := fs.generateObjectKey(info.VirtualPath, &info.Name)

	if err := fs.handler.Put(ctx, info.File, objectKey, info.Size); err != nil {
		return err
	}

	return fs.Trigger(ctx, HookAfterUpload)
}

func (fs *FileSystem) Download(ctx context.Context, info DownloadFileInfo) (io.ReadSeekCloser, error) {
	objectKey := fs.generateObjectKey(info.VirtualPath, info.Name)

	return fs.handler.Get(ctx, objectKey)
}

func (fs *FileSystem) generateObjectKey(virtualPath string, name *string) string {
	folder := fmt.Sprintf("uploads/%d/%s", fs.User.ID, virtualPath)
	if name == nil {
		return folder
	} else {
		return path.Join(folder, *name)
	}
}

package filesystem

import (
	"context"
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

type FileInfo struct {
	File     io.Reader
	Name     string
	Size     uint64
	MIMEType string
	Path     string
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

func (fs *FileSystem) Upload(ctx context.Context, info FileInfo) error {
	err := fs.handler.Put(ctx, info.File, "", info.Size)
	if err != nil {
		return err
	}

	err = fs.Trigger(ctx, HookAfterUpload)
	if err != nil {
		return err
	}

	return nil
}

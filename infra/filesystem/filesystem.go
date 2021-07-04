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

type FileInfo struct {
	File        io.Reader
	Name        string
	Size        uint64
	MIMEType    string
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

func (fs *FileSystem) Upload(ctx context.Context, info FileInfo) error {
	objectKey := fs.generateObjectKey(info)
	err := fs.handler.Put(ctx, info.File, objectKey, info.Size)
	if err != nil {
		return err
	}

	err = fs.Trigger(HookAfterUpload, info)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) generateObjectKey(info FileInfo) string {
	return path.Join(fmt.Sprintf("uploads/%d/%s", fs.User.ID, info.VirtualPath), info.Name)
}

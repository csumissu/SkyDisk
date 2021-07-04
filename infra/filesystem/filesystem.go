package filesystem

import (
	"github.com/csumissu/SkyDisk/infra/filesystem/handlers/local"
	"github.com/csumissu/SkyDisk/models"
	"sync"
)

type FileSystem struct {
	User    *models.User
	Handler Handler
	Hooks   map[string][]Hook
	mutex   sync.Mutex
}

func NewFileSystem(user *models.User) (*FileSystem, error) {
	fs := getEmptyFS()
	fs.User = user
	fs.determineHandler()
	return fs, nil
}

func (fs *FileSystem) determineHandler() {
	fs.Handler = local.Handler{}
}

package filesystem

import "sync"

var fsPool = sync.Pool{
	New: func() interface{} {
		return &FileSystem{}
	},
}

func getEmptyFS() *FileSystem {
	fs := fsPool.Get().(*FileSystem)
	return fs
}

func (fs *FileSystem) Recycle() {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fs.reset()
	fsPool.Put(fs)
}

func (fs *FileSystem) reset() {
	fs.User = nil
	fs.handler = nil
	fs.hooks = nil
}

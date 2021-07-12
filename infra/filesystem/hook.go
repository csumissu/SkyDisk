package filesystem

import (
	"context"
	"github.com/csumissu/SkyDisk/util"
)

const HookAfterUploadFile = "AfterUploadFile"
const HookAfterDeleteObject = "AfterDeleteObject"
const HookAfterCreateDir = "AfterCreateDir"

type Hook func(ctx context.Context, fs *FileSystem) error

func (fs *FileSystem) Use(name string, hook Hook) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if fs.hooks == nil {
		fs.hooks = make(map[string][]Hook)
	}
	if _, ok := fs.hooks[name]; ok {
		fs.hooks[name] = append(fs.hooks[name], hook)
	} else {
		fs.hooks[name] = []Hook{hook}
	}
}

func (fs *FileSystem) Trigger(ctx context.Context, name string) error {
	fs.mutex.RLocker().Lock()
	defer fs.mutex.RLocker().Unlock()

	if hooks, ok := fs.hooks[name]; ok {
		for _, hook := range hooks {
			if err := hook(ctx, fs); err != nil {
				util.Logger.Warn("trigger hook failed, name: %s, hook: %v", name, hook, err)
				return err
			}
		}
	}
	return nil
}

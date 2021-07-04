package filesystem

import (
	"context"
	"github.com/csumissu/SkyDisk/util"
)

type Hook func(ctx context.Context, fs *FileSystem) error

func (fs *FileSystem) Use(name string, hook Hook) {
	if _, ok := fs.Hooks[name]; ok {
		fs.Hooks[name] = append(fs.Hooks[name], hook)
	} else {
		fs.Hooks[name] = []Hook{hook}
	}
}

func (fs *FileSystem) Trigger(ctx context.Context, name string) error {
	if hooks, ok := fs.Hooks[name]; ok {
		for _, hook := range hooks {
			if err := hook(ctx, fs); err != nil {
				util.Logger.Warn("trigger hook failed, name: %s, hook: %v", name, hook, err)
				return err
			}
		}
	}
	return nil
}

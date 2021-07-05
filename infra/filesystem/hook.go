package filesystem

import (
	"github.com/csumissu/SkyDisk/util"
)

const HookAfterUpload = "AfterUpload"

type HookParams map[string]interface{}

type Hook func(fs *FileSystem, params HookParams) error

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

func (fs *FileSystem) Trigger(name string, params ...interface{}) error {
	fs.mutex.RLocker().Lock()
	defer fs.mutex.RLocker().Unlock()

	if hooks, ok := fs.hooks[name]; ok {
		actualParams := util.Convert(params...)
		for _, hook := range hooks {
			if err := hook(fs, actualParams); err != nil {
				util.Logger.Warn("trigger hook failed, name: %s, hook: %v", name, hook, err)
				return err
			}
		}
	}
	return nil
}

func (p HookParams) MustGet(i interface{}) interface{} {
	return p[util.GetTypeName(i)]
}

package local

import (
	"context"
	"github.com/csumissu/SkyDisk/util"
	"io"
	"os"
	"path/filepath"
)

type Handler struct {
}

var rootDir = getRootDir()

func (handler Handler) Put(ctx context.Context, file io.Reader, objectKey string, size uint64) error {
	savePath := filepath.Join(rootDir, objectKey)
	util.Logger.Debug("uploading file to %s, size: %d", savePath, size)

	basePath := filepath.Dir(savePath)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		err := os.MkdirAll(basePath, 0744)
		if err != nil {
			util.Logger.Warn("directory could not be created, basePath: %s", basePath, err)
			return err
		}
	}

	target, err := os.Create(savePath)
	if err != nil {
		util.Logger.Warn("file could not be created, objectKey: %s", savePath, err)
		return err
	}

	_, err = io.Copy(target, file)
	return err
}

func getRootDir() string {
	ex, err := os.Executable()
	if err != nil {
		util.Logger.Panic("executable path could not be retrieved", err)
	}
	return filepath.Dir(ex)
}

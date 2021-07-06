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
	fullPath := filepath.Join(rootDir, objectKey)
	util.Logger.Debug("uploading file to %s, size: %d", fullPath, size)

	basePath := filepath.Dir(fullPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.MkdirAll(basePath, 0744)
		if err != nil {
			util.Logger.Warn("directory could not be created, basePath: %s", basePath, err)
			return err
		}
	}

	target, err := os.Create(fullPath)
	if err != nil {
		util.Logger.Warn("file could not be created, fullPath: %s", fullPath, err)
		return err
	}

	_, err = io.Copy(target, file)
	return err
}

func (handler Handler) Get(ctx context.Context, objectKey string) (io.ReadSeekCloser, error) {
	fullPath := filepath.Join(rootDir, objectKey)
	if file, err := os.Open(fullPath); err != nil {
		util.Logger.Warn("file could not be opened, fullPath: %s", fullPath, err)
		return nil, err
	} else {
		return file, nil
	}
}

func getRootDir() string {
	ex, err := os.Executable()
	if err != nil {
		util.Logger.Panic("executable path could not be retrieved", err)
	}
	return filepath.Dir(ex)
}

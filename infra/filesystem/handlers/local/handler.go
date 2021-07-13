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

const defaultDirPermission = 0744

func (handler Handler) Put(ctx context.Context, file io.Reader, objectKey string, size uint64) error {
	fullPath := filepath.Join(rootDir, objectKey)
	util.Logger.Debug("uploading file, fullPath: %v, size: %v", fullPath, size)

	basePath := filepath.Dir(fullPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err = os.MkdirAll(basePath, defaultDirPermission); err != nil {
			util.Logger.Warn("directory could not be created, basePath: %v", basePath, err)
			return err
		}
	}

	target, err := os.Create(fullPath)
	if err != nil {
		util.Logger.Warn("file could not be created, fullPath: %v", fullPath, err)
		return err
	}

	_, err = io.Copy(target, file)
	return err
}

func (handler Handler) Get(ctx context.Context, objectKey string, isDir bool) (io.ReadSeekCloser, error) {
	fullPath := filepath.Join(rootDir, objectKey)
	util.Logger.Debug("downloading object, fullPath: %v", fullPath)

	if file, err := os.Open(fullPath); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func (handler Handler) Delete(ctx context.Context, objectKey string, isDir bool) error {
	fullPath := filepath.Join(rootDir, objectKey)
	util.Logger.Debug("deleting object, fullPath: %v", fullPath)

	var err error
	if isDir {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}

	if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}
}

func (handler Handler) CreateDir(ctx context.Context, objectKey string) error {
	fullPath := filepath.Join(rootDir, objectKey)
	util.Logger.Debug("creating dir, fullPath: %v", fullPath)

	return os.MkdirAll(fullPath, defaultDirPermission)
}

func (handler Handler) Rename(ctx context.Context, srcObjectKey string, dstObjectKey string) error {
	srcFullPath := filepath.Join(rootDir, srcObjectKey)
	dstFullPath := filepath.Join(rootDir, dstObjectKey)
	util.Logger.Debug("renaming object, src: %v, dst: %v", srcFullPath, dstFullPath)

	if err := os.MkdirAll(filepath.Dir(dstFullPath), defaultDirPermission); err != nil {
		return err
	}

	return os.Rename(srcFullPath, dstFullPath)
}

func (handler Handler) Move(ctx context.Context, srcObjectKey string, dstObjectKey string) error {
	srcFullPath := filepath.Join(rootDir, srcObjectKey)
	dstFullPath := filepath.Join(rootDir, dstObjectKey)
	util.Logger.Debug("moving object, src: %v, dst: %v", srcFullPath, dstFullPath)

	if err := os.MkdirAll(filepath.Dir(dstFullPath), defaultDirPermission); err != nil {
		return err
	}

	return os.Rename(srcFullPath, dstFullPath)
}

func getRootDir() string {
	ex, err := os.Executable()
	if err != nil {
		util.Logger.Panic("executable path could not be retrieved", err)
	}
	return filepath.Dir(ex)
}

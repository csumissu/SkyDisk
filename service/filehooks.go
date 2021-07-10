package service

import (
	"context"
	"github.com/csumissu/SkyDisk/infra/filesystem"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/util"
	"gorm.io/gorm"
	"path"
	"strings"
)

func HookGenericAfterUpload(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.UploadObjectInfoCtx).(filesystem.UploadObjectInfo)
	util.Logger.Debug("hook genericAfterUpload, fileInfo: %v", info)

	dir, err := createDirsRecursively(*fs.User, path.Dir(info.VirtualPath))
	if err != nil {
		return err
	}

	return createOrUpdateFile(*fs.User, info, *dir)
}

func HookGenericAfterDelete(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.DeleteObjectInfoCtx).(filesystem.DeleteObjectInfo)
	util.Logger.Debug("hook genericAfterDelete, objectInfo: %v", info)

	if object, err := fs.User.GetObjectByID(info.ObjectID); err == nil {
		return object.Delete()
	} else if err == gorm.ErrRecordNotFound {
		return nil
	} else {
		return err
	}
}

func createOrUpdateFile(user models.User, info filesystem.UploadObjectInfo, dir models.Object) error {
	file, err := user.GetFileByNameAndDirID(info.Name, dir.ID)
	if err == nil {
		file.MIMEType = &info.MIMEType
		file.Size = &info.Size
		return file.Update()
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	_, err = user.NewFile(dir, info.Name, info.Size, info.MIMEType)
	return err
}

func createDirsRecursively(user models.User, fullPath string) (*models.Object, error) {
	fullPath = path.Clean(fullPath)
	parentDirPath := path.Dir(fullPath)
	currentDirName := strings.TrimRight(path.Base(fullPath), " ")

	if fullPath == "/" || fullPath == "." || fullPath == "" {
		rootDir, err := user.GetRootDir()
		if err == nil {
			return rootDir, nil
		} else if err == gorm.ErrRecordNotFound {
			return user.CreateRootDir()
		} else {
			return nil, err
		}
	}

	dir, err := user.GetDirByFullPath(fullPath)
	if err == nil {
		return dir, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	parentDir, err := createDirsRecursively(user, parentDirPath)
	if err != nil {
		return nil, err
	}

	return user.NewFolder(*parentDir, currentDirName)
}

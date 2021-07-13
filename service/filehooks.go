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

func HookGenericAfterUploadFile(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.UploadFileInfoCtx).(filesystem.UploadFileInfo)
	util.Logger.Debug("hook genericAfterUploadFile, info: %v", info)

	dir, err := createDirsRecursively(*fs.User, path.Dir(info.VirtualPath))
	if err != nil {
		return err
	}

	return createOrUpdateFile(*fs.User, info, *dir)
}

func HookGenericAfterDeleteObject(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.DeleteObjectInfoCtx).(filesystem.DeleteObjectInfo)
	util.Logger.Debug("hook genericAfterDeleteObject, info: %v", info)

	if object, err := fs.User.GetObjectByID(info.ObjectID); err == nil {
		return object.Delete()
	} else if err == gorm.ErrRecordNotFound {
		return nil
	} else {
		return err
	}
}

func HookGenericAfterCreateDir(ctx context.Context, fs *filesystem.FileSystem) error {
	virtualPath := ctx.Value(filesystem.CreateDirCtx).(string)
	util.Logger.Debug("hook genericAfterCreateDir, path: %v", virtualPath)

	_, err := createDirsRecursively(*fs.User, virtualPath)
	return err
}

func HookGenericAfterRenameObject(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.RenameObjectInfoCtx).(filesystem.RenameObjectInfo)
	util.Logger.Debug("hook genericAfterRenameObject, info: %v", info)

	if object, err := fs.User.GetObjectByID(info.ObjectID); err == nil {
		return object.Rename(info.DstVirtualPath)
	} else {
		return nil
	}
}

func HookGenericAfterMoveObject(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.MoveObjectInfoCtx).(filesystem.MoveObjectInfo)
	util.Logger.Debug("hook genericAfterMoveObject, info: %v", info)

	dir, err := createDirsRecursively(*fs.User, path.Dir(info.DstVirtualPath))
	if err != nil {
		return err
	}

	if object, err := fs.User.GetObjectByID(info.ObjectID); err == nil {
		return object.MoveTo(*dir)
	} else {
		return nil
	}
}

func createOrUpdateFile(user models.User, info filesystem.UploadFileInfo, dir models.Object) error {
	file, err := user.GetObjectByNameAndParentID(info.Name, dir.ID, models.FILE)
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

	currentDir, err := user.GetObjectByFullPath(fullPath, models.DIR)
	if err == nil {
		return currentDir, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	parentDir, err := createDirsRecursively(user, parentDirPath)
	if err != nil {
		return nil, err
	}

	return user.NewDir(*parentDir, currentDirName)
}

package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/infra/filesystem"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/util"
	"gorm.io/gorm"
	"path"
	"strings"
)

type FileService struct {
}

func (service *FileService) SingleUpload(ctx context.Context, path string, upload graphql.Upload) (bool, error) {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return false, err
	}

	fs, err := filesystem.NewFileSystem(user)
	if err != nil {
		return false, err
	}
	defer fs.Recycle()

	fileInfo := filesystem.FileInfo{
		File:        upload.File,
		Name:        upload.Filename,
		Size:        uint64(upload.Size),
		MIMEType:    upload.ContentType,
		VirtualPath: path,
	}

	fs.Use(filesystem.HookAfterUpload, HookGenericAfterUpload)
	err = fs.Upload(ctx, fileInfo)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service *FileService) ListObjects(ctx context.Context, path string) ([]*dto.ObjectResponse, error) {
	return make([]*dto.ObjectResponse, 0), nil
}

func HookGenericAfterUpload(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.FileInfoCtx).(filesystem.FileInfo)
	util.Logger.Debug("genericAfterUpload, fileInfo: %v", info)

	folder, err := createDirectory(fs.User.ID, info.VirtualPath)
	if err != nil {
		return err
	}

	return createOrUpdateFile(fs.User.ID, info, folder)
}

func createOrUpdateFile(userID uint, info filesystem.FileInfo, folder *models.Folder) error {
	file, err := models.GetFileByNameAndFolderID(userID, info.Name, folder.ID)
	if err == nil {
		file.MIMEType = info.MIMEType
		file.Size = info.Size
		return file.Update()
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	file = &models.File{
		Name:     info.Name,
		UserID:   userID,
		Size:     info.Size,
		MIMEType: info.MIMEType,
		FolderID: folder.ID,
	}
	return file.Create()
}

func createDirectory(userID uint, fullPath string) (*models.Folder, error) {
	fullPath = path.Clean(fullPath)
	parentDir := path.Dir(fullPath)
	currentDir := path.Base(fullPath)
	currentDir = strings.TrimRight(currentDir, " ")

	if fullPath == "/" || fullPath == "." || fullPath == "" {
		rootFolder, err := models.GetRootFolder(userID)
		if err == nil {
			return rootFolder, nil
		} else if err == gorm.ErrRecordNotFound {
			return models.CreateRootFolder(userID)
		} else {
			return nil, err
		}
	}

	currentFolder, err := models.GetFolderByFullPath(userID, fullPath)
	if err == nil {
		return currentFolder, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	parentFolder, err := createDirectory(userID, parentDir)
	if err != nil {
		return nil, err
	}

	currentFolder = &models.Folder{
		Name:     currentDir,
		FullPath: path.Join(parentFolder.FullPath, currentDir),
		ParentID: &parentFolder.ID,
		UserID:   userID,
	}
	err = currentFolder.Create()
	if err != nil {
		return nil, err
	}

	return currentFolder, nil
}

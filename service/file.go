package service

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/infra/filesystem"
	"github.com/csumissu/SkyDisk/models"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type FileService struct {
}

const (
	DIR  string = "dir"
	FILE string = "file"
)

func (service *FileService) SingleUpload(ctx context.Context, virtualPath string, upload graphql.Upload) (bool, error) {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return false, err
	}

	fs, err := filesystem.NewFileSystem(user)
	if err != nil {
		return false, err
	}
	defer fs.Recycle()

	fileInfo := filesystem.UploadFileInfo{
		File:        upload.File,
		Name:        upload.Filename,
		Size:        uint64(upload.Size),
		MIMEType:    upload.ContentType,
		VirtualPath: path.Clean(virtualPath),
	}

	fs.Use(filesystem.HookAfterUpload, HookGenericAfterUpload)
	err = fs.Upload(ctx, fileInfo)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service *FileService) ListObjects(ctx context.Context, virtualPath string) ([]*dto.ObjectResponse, error) {
	virtualPath = path.Clean(virtualPath)
	userID := util.MustGetCurrentUserID(ctx)

	folder, err := models.GetFolderByFullPath(userID, virtualPath)
	if err == gorm.ErrRecordNotFound {
		return make([]*dto.ObjectResponse, 0), nil
	} else if err != nil {
		return nil, err
	}

	childFolders, err := folder.GetChildFolders()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	childFiles, err := folder.GetChildFiles()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	objects := make([]*dto.ObjectResponse, 0, len(childFolders)+len(childFiles))
	for _, childFolder := range childFolders {
		objects = append(objects, &dto.ObjectResponse{
			ID:        childFolder.ID,
			Name:      childFolder.Name,
			Path:      childFolder.FullPath,
			Type:      DIR,
			MimeType:  nil,
			UpdatedAt: childFolder.UpdatedAt,
			CreatedAt: childFolder.CreatedAt,
		})
	}
	for _, childFile := range childFiles {
		objects = append(objects, &dto.ObjectResponse{
			ID:        childFile.ID,
			Name:      childFile.Name,
			Path:      path.Join(folder.FullPath, childFile.Name),
			Type:      FILE,
			MimeType:  &childFile.MIMEType,
			UpdatedAt: childFile.UpdatedAt,
			CreatedAt: childFile.CreatedAt,
		})
	}

	return objects, nil
}

func (service *FileService) Download(c *gin.Context, objectId uint) dto.Response {
	ctx := c.Request.Context()
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return dto.Failure(http.StatusUnauthorized, err.Error())
	}

	fs, err := filesystem.NewFileSystem(user)
	if err != nil {
		return dto.Failure(http.StatusInternalServerError, err.Error())
	}
	defer fs.Recycle()

	file, err := models.GetFileByID(user.ID, objectId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return dto.Failure(http.StatusNotFound, err.Error())
	} else if err != nil {
		return dto.Failure(http.StatusInternalServerError, err.Error())
	}

	folder, err := models.GetFolderByID(user.ID, file.FolderID)
	if err != nil {
		return dto.Failure(http.StatusInternalServerError, err.Error())
	}

	fileInfo := filesystem.DownloadFileInfo{
		Name:        &file.Name,
		VirtualPath: folder.FullPath,
	}
	object, err := fs.Download(ctx, fileInfo)
	if err != nil {
		return dto.Failure(http.StatusInternalServerError, err.Error())
	}
	defer func(object io.ReadSeekCloser) {
		_ = object.Close()
	}(object)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.PathEscape(file.Name)))
	http.ServeContent(c.Writer, c.Request, file.Name, file.UpdatedAt, object)

	return dto.EmptyResponse()
}

func HookGenericAfterUpload(ctx context.Context, fs *filesystem.FileSystem) error {
	info := ctx.Value(filesystem.UploadFileInfoCtx).(filesystem.UploadFileInfo)
	util.Logger.Debug("genericAfterUpload, fileInfo: %v", info)

	folder, err := createDirectory(fs.User.ID, info.VirtualPath)
	if err != nil {
		return err
	}

	return createOrUpdateFile(fs.User.ID, info, folder)
}

func createOrUpdateFile(userID uint, info filesystem.UploadFileInfo, folder *models.Folder) error {
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

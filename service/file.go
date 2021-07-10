package service

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/csumissu/SkyDisk/infra/filesystem"
	"github.com/csumissu/SkyDisk/routers/dto"
	"github.com/csumissu/SkyDisk/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type FileService struct {
}

func (service *FileService) SingleUpload(ctx context.Context, virtualPath string, upload graphql.Upload) (bool, error) {
	fs, err := getFileSystem(ctx)
	if err != nil {
		return false, err
	}
	defer fs.Recycle()

	virtualPath = path.Clean(virtualPath)
	if !strings.HasSuffix(virtualPath, upload.Filename) {
		virtualPath = path.Join(virtualPath, upload.Filename)
	}

	fileInfo := filesystem.UploadObjectInfo{
		Name:        upload.Filename,
		Size:        uint64(upload.Size),
		MIMEType:    upload.ContentType,
		VirtualPath: virtualPath,
	}

	fs.Use(filesystem.HookAfterUpload, HookGenericAfterUpload)
	err = fs.Upload(ctx, upload.File, fileInfo)
	return err == nil, err
}

func (service *FileService) ListObjects(ctx context.Context, virtualPath string) ([]*dto.ObjectResponse, error) {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	virtualPath = path.Clean(virtualPath)
	dir, err := user.GetDirByFullPath(virtualPath)
	if err == gorm.ErrRecordNotFound {
		return make([]*dto.ObjectResponse, 0), nil
	} else if err != nil {
		return nil, err
	}

	childObjects, err := dir.GetChildObjects()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	results := make([]*dto.ObjectResponse, 0, len(childObjects))
	for _, object := range childObjects {
		results = append(results, &dto.ObjectResponse{
			ID:        object.ID,
			Name:      object.Name,
			Path:      object.FullPath,
			Type:      object.GetType(),
			Size:      object.Size,
			MimeType:  object.MIMEType,
			UpdatedAt: object.UpdatedAt,
			CreatedAt: object.CreatedAt,
		})
	}

	return results, nil
}

func (service *FileService) DownloadObject(c *gin.Context, objectID uint) dto.Response {
	ctx := c.Request.Context()
	fs, err := getFileSystem(ctx)
	if err != nil {
		return dto.Failure(http.StatusUnauthorized, err.Error())
	}
	defer fs.Recycle()

	object, err := fs.User.GetObjectByID(objectID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return dto.Failure(http.StatusNotFound, err.Error())
	} else if err != nil {
		return dto.Failure(http.StatusInternalServerError, err.Error())
	}

	objectInfo := filesystem.DownloadObjectInfo{
		IsDir:       object.IsDir(),
		VirtualPath: object.FullPath,
	}

	content, err := fs.Download(ctx, objectInfo)
	if err != nil {
		return dto.Failure(http.StatusInternalServerError, err.Error())
	}
	defer util.CloseQuietly(content)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.PathEscape(object.Name)))
	http.ServeContent(c.Writer, c.Request, object.Name, object.UpdatedAt, content)

	return dto.EmptyResponse()
}

func (service *FileService) DeleteObject(ctx context.Context, objectID uint) (bool, error) {
	fs, err := getFileSystem(ctx)
	if err != nil {
		return false, err
	}
	defer fs.Recycle()

	object, err := fs.User.GetObjectByID(objectID)
	if err != nil {
		return false, err
	}

	info := filesystem.DeleteObjectInfo{
		ObjectID:    object.ID,
		IsDir:       object.IsDir(),
		VirtualPath: object.FullPath,
	}

	fs.Use(filesystem.HookAfterDelete, HookGenericAfterDelete)
	err = fs.Delete(ctx, info)
	return err == nil, err
}

func getFileSystem(ctx context.Context) (*filesystem.FileSystem, error) {
	if user, err := GetCurrentUser(ctx); err == nil {
		return filesystem.NewFileSystem(user)
	} else {
		return nil, err
	}
}

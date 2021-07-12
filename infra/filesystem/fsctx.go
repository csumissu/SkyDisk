package filesystem

type key int

const (
	UploadFileInfoCtx key = iota
	DownloadObjectInfoCtx
	DeleteObjectInfoCtx
	CreateDirCtx
)

type UploadFileInfo struct {
	Name        string
	Size        uint64
	MIMEType    string
	VirtualPath string
}

type DownloadObjectInfo struct {
	ObjectID    uint
	IsDir       bool
	VirtualPath string
}

type DeleteObjectInfo struct {
	ObjectID    uint
	IsDir       bool
	VirtualPath string
}

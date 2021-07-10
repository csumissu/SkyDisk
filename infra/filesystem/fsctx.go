package filesystem

type key int

const (
	UploadObjectInfoCtx key = iota
	DownloadObjectInfoCtx
	DeleteObjectInfoCtx
)

type UploadObjectInfo struct {
	Name        string
	Size        uint64
	MIMEType    string
	VirtualPath string
}

type DownloadObjectInfo struct {
	IsDir       bool
	VirtualPath string
}

type DeleteObjectInfo struct {
	ObjectID    uint
	IsDir       bool
	VirtualPath string
}

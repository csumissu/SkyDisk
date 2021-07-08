package filesystem

import "io"

type key int

const (
	UploadFileInfoCtx key = iota
	DeleteObjectInfoCtx
)

type UploadFileInfo struct {
	File        io.Reader
	Name        string
	Size        uint64
	MIMEType    string
	VirtualPath string
}

type DownloadObjectInfo struct {
	Name        *string
	VirtualPath string
}

type DeleteObjectInfo struct {
	ObjectID    uint
	Name        *string
	VirtualPath string
}

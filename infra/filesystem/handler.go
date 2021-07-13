package filesystem

import (
	"context"
	"io"
)

type Handler interface {
	Put(ctx context.Context, file io.Reader, objectKey string, size uint64) error

	Get(ctx context.Context, objectKey string, isDir bool) (io.ReadSeekCloser, error)

	Delete(ctx context.Context, objectKey string, isDir bool) error

	CreateDir(ctx context.Context, objectKey string) error

	Rename(ctx context.Context, srcObjectKey string, dstObjectKey string) error

	Move(ctx context.Context, srcObjectKey string, dstObjectKey string) error
}

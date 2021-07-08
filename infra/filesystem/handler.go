package filesystem

import (
	"context"
	"io"
)

type Handler interface {
	Put(ctx context.Context, file io.Reader, objectKey string, size uint64) error

	Get(ctx context.Context, objectKey string) (io.ReadSeekCloser, error)

	Delete(ctx context.Context, objectKey string) error
}

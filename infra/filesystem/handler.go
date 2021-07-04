package filesystem

import (
	"context"
	"io"
)

type Handler interface {
	Put(ctx context.Context, file io.Reader, objectKey string, size uint64) error
}

package filesystem

import (
	"context"
	"io"
)

type Handler interface {
	Put(ctx context.Context, file io.ReadCloser, dest string, size uint64) error
}

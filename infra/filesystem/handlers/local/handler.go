package local

import (
	"context"
	"io"
)

type Handler struct {
}

func (handler Handler) Put(ctx context.Context, file io.ReadCloser, dest string, size uint64) error {
	return nil
}

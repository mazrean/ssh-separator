package workspace

import (
	"context"
	"errors"
	"io"
)

var (
	ErrWorkspaceExist = errors.New("workspace exist error")
)

type IWorkspace interface {
	Create(ctx context.Context, userName string) error
	Connect(ctx context.Context, userName string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error
	Remove(ctx context.Context, userName string) error
}

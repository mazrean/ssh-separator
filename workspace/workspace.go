//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package workspace

import (
	"context"
	"errors"
	"io"

	"github.com/mazrean/separated-webshell/domain"
)

var (
	ErrWorkspaceExist = errors.New("workspace exist error")
)

type IWorkspace interface {
	Create(ctx context.Context, userName string) error
	Connect(ctx context.Context, userName string, isTty bool, winCh <-chan *domain.Window, stdin io.Reader, stdout io.Writer, stderr io.Writer) error
	Remove(ctx context.Context, userName string) error
}

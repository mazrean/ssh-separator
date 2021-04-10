//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package workspace

import (
	"context"
	"errors"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

var (
	ErrWorkspaceExist = errors.New("workspace exist error")
)

type IWorkspace interface {
	Create(ctx context.Context, userName values.UserName) (*domain.Workspace, error)
	Start(ctx context.Context, workspace *domain.Workspace) error
	Stop(ctx context.Context, workspace *domain.Workspace) error
}

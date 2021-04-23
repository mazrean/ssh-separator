//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package workspace

import (
	"context"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

type IWorkspaceConnection interface {
	Connect(ctx context.Context, workspace *domain.Workspace) (*domain.WorkspaceConnection, error)
	Disconnect(ctx context.Context, connection *domain.WorkspaceConnection) error
	Resize(ctx context.Context, connection *domain.WorkspaceConnection, window *values.Window) error
}

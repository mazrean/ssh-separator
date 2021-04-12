//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package store

import (
	"context"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

type IWorkspace interface {
	Set(ctx context.Context, userName values.UserName, workspace *domain.Workspace) error
	Get(ctx context.Context, userName values.UserName) (*domain.Workspace, error)
}

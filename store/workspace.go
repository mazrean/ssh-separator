//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package store

import (
	"context"
	"errors"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

var (
	// ErrWorkspaceNotFound a workspace is not found.
	ErrWorkspaceNotFound = errors.New("not found")
)

type IWorkspace interface {
	Set(ctx context.Context, userName values.UserName, workspace *domain.Workspace) error
	Get(ctx context.Context, userName values.UserName) (*domain.Workspace, error)
}

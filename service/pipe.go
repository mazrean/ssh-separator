package service

import (
	"context"
	"fmt"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/workspace"
)

type IPipe interface {
	Pipe(ctx context.Context, userName values.UserName, connection *domain.Connection) error
}

type Pipe struct {
	workspace.IWorkspace
}

func NewPipe(w workspace.IWorkspace) *Pipe {
	return &Pipe{
		IWorkspace: w,
	}
}

func (p *Pipe) Pipe(ctx context.Context, userName values.UserName, connection *domain.Connection) error {
	err := p.IWorkspace.Connect(ctx, userName, connection)
	if err != nil {
		return fmt.Errorf("connect to workspace error: %w", err)
	}

	return nil
}

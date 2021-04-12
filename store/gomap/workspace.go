package gomap

import (
	"context"
	"errors"
	"sync"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

type Workspace struct {
	sync.Map
}

func NewWorkspace() *Workspace {
	return &Workspace{
		Map: sync.Map{},
	}
}

func (w *Workspace) Set(ctx context.Context, userName values.UserName, workspace *domain.Workspace) error {
	w.Map.Store(userName, workspace)

	return nil
}

func (w *Workspace) Get(ctx context.Context, userName values.UserName) (*domain.Workspace, error) {
	iWorkspace, ok := w.Map.Load(userName)
	if !ok {
		return nil, errors.New("user not found")
	}

	workspace, ok := iWorkspace.(*domain.Workspace)
	if !ok {
		return nil, errors.New("workspace is broken")
	}

	return workspace, nil
}

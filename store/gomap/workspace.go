package gomap

import (
	"context"
	"errors"
	"sync"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/store"
)

type Workspace struct {
	syncMap sync.Map
}

func NewWorkspace() *Workspace {
	return &Workspace{
		syncMap: sync.Map{},
	}
}

func (w *Workspace) Set(ctx context.Context, userName values.UserName, workspace *domain.Workspace) error {
	w.syncMap.Store(userName, workspace)

	return nil
}

func (w *Workspace) Get(ctx context.Context, userName values.UserName) (*domain.Workspace, error) {
	iWorkspace, ok := w.syncMap.Load(userName)
	if !ok {
		return nil, store.ErrWorkspaceNotFound
	}

	workspace, ok := iWorkspace.(*domain.Workspace)
	if !ok {
		return nil, errors.New("workspace is broken")
	}

	return workspace, nil
}

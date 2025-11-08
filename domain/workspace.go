package domain

import (
	"errors"
	"sync/atomic"

	"github.com/mazrean/separated-webshell/domain/values"
)

type Workspace struct {
	id             values.WorkspaceID
	name           values.WorkspaceName
	userName       values.UserName
	Status         values.WorkspaceStatus
	connectionNum  int64
	maxConnections int64
}

func NewWorkspace(id values.WorkspaceID, name values.WorkspaceName, userName values.UserName, maxConnections int64) *Workspace {
	return &Workspace{
		id:             id,
		name:           name,
		userName:       userName,
		Status:         values.StatusDown,
		connectionNum:  0,
		maxConnections: maxConnections,
	}
}

func (w *Workspace) ID() values.WorkspaceID {
	return w.id
}

func (w *Workspace) Name() values.WorkspaceName {
	return w.name
}

func (w *Workspace) UserName() values.UserName {
	return w.userName
}

func (w *Workspace) ConnectionNum() int64 {
	return w.connectionNum
}

func (w *Workspace) AddConnection() error {
	// Check if adding a connection would exceed the limit
	for {
		current := atomic.LoadInt64(&w.connectionNum)
		if current >= w.maxConnections {
			return ErrTooManyConnections
		}
		if atomic.CompareAndSwapInt64(&w.connectionNum, current, current+1) {
			return nil
		}
	}
}

func (w *Workspace) RemoveConnection() error {
	if w.connectionNum <= 0 {
		return errors.New("no connection")
	}

	atomic.AddInt64(&w.connectionNum, -1)

	return nil
}

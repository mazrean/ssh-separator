package domain

import (
	"errors"
	"sync/atomic"

	"github.com/mazrean/separated-webshell/domain/values"
)

type Workspace struct {
	id            values.WorkspaceID
	name          values.WorkspaceName
	userName      values.UserName
	Status        values.WorkspaceStatus
	connectionNum int32
}

func NewWorkspace(id values.WorkspaceID, name values.WorkspaceName, userName values.UserName) *Workspace {
	return &Workspace{
		id:            id,
		name:          name,
		userName:      userName,
		Status:        values.StatusDown,
		connectionNum: 0,
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

func (w *Workspace) ConnectionNum() int32 {
	return w.connectionNum
}

func (w *Workspace) AddConnection() error {
	atomic.AddInt32(&w.connectionNum, 1)

	return nil
}

func (w *Workspace) RemoveConnection() error {
	if w.connectionNum <= 0 {
		return errors.New("no connection")
	}

	atomic.AddInt32(&w.connectionNum, -1)

	return nil
}

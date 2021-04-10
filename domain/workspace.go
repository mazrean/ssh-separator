package domain

import "github.com/mazrean/separated-webshell/domain/values"

type Workspace struct {
	name     values.WorkspaceName
	userName values.UserName
}

func NewWorkspace(name values.WorkspaceName, userName values.UserName) *Workspace {
	return &Workspace{
		name:     name,
		userName: userName,
	}
}

func (w *Workspace) Name() values.WorkspaceName {
	return w.name
}

func (w *Workspace) UserName() values.UserName {
	return w.userName
}

package domain

import (
	"io"

	"github.com/mazrean/separated-webshell/domain/values"
)

type WorkspaceConnection struct {
	id values.WorkspaceConnectionID
	io *values.WorkspaceIO
}

func NewWorkspaceConnection(id values.WorkspaceConnectionID, io *values.WorkspaceIO) *WorkspaceConnection {
	return &WorkspaceConnection{
		id: id,
		io: io,
	}
}

func (wc *WorkspaceConnection) ID() values.WorkspaceConnectionID {
	return wc.id
}

func (wc *WorkspaceConnection) WriteCloser() io.WriteCloser {
	return wc.io.WriteCloser()
}

func (wc *WorkspaceConnection) ReadCloser() io.ReadCloser {
	return wc.io.ReadCloser()
}

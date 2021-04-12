package values

import "io"

type (
	WorkspaceConnectionID string
)

func NewWorkspaceConnectionID(id string) WorkspaceConnectionID {
	return WorkspaceConnectionID(id)
}

type WorkspaceIO struct {
	writer io.WriteCloser
	reader io.ReadCloser
}

func NewWorkspaceIO(writer io.WriteCloser, reader io.ReadCloser) *WorkspaceIO {
	return &WorkspaceIO{
		writer: writer,
		reader: reader,
	}
}

func (wio *WorkspaceIO) WriteCloser() io.WriteCloser {
	return wio.writer
}

func (wio *WorkspaceIO) ReadCloser() io.ReadCloser {
	return wio.reader
}

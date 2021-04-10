package values

import "io"

type (
	WorkspaceConnectionID string
)

func NewWorkspaceConnectionID(id string) WorkspaceConnectionID {
	return WorkspaceConnectionID(id)
}

type WorkspaceIO struct {
	writer io.Writer
	reader io.Reader
}

func NewWorkspaceIO(writer io.Writer, reader io.Reader) *WorkspaceIO {
	return &WorkspaceIO{
		writer: writer,
		reader: reader,
	}
}

func (wio *WorkspaceIO) Writer() io.Writer {
	return wio.writer
}

func (wio *WorkspaceIO) Reader() io.Reader {
	return wio.reader
}

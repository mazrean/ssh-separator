package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

var (
	createOpts = types.ExecConfig{
		User:         imageUser,
		WorkingDir:   fmt.Sprintf("/home/%s", imageUser),
		Cmd:          []string{imageCmd},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}
	attachOpts = types.ExecStartCheck{
		Tty: true,
	}
)

type WorkspaceConnection struct{}

func NewWorkspaceConnection() *WorkspaceConnection {
	return &WorkspaceConnection{}
}

func (wc *WorkspaceConnection) Connect(ctx context.Context, workspace *domain.Workspace) (*domain.WorkspaceConnection, error) {
	idRes, err := cli.ContainerExecCreate(ctx, string(workspace.ID()), createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec: %w", err)
	}

	stream, err := cli.ContainerExecAttach(ctx, idRes.ID, attachOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to attach container: %w", err)
	}

	connectionID := values.NewWorkspaceConnectionID(idRes.ID)
	connectionIO := values.NewWorkspaceIO(stream.Conn, io.NopCloser(stream.Reader))

	return domain.NewWorkspaceConnection(connectionID, connectionIO), nil
}

func (wc *WorkspaceConnection) Disconnect(ctx context.Context, connection *domain.WorkspaceConnection) error {
	err := connection.ReadCloser().Close()
	if err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}

	err = connection.WriteCloser().Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func (wc *WorkspaceConnection) Resize(ctx context.Context, connection *domain.WorkspaceConnection, window *values.Window) error {
	err := cli.ContainerExecResize(ctx, string(connection.ID()), types.ResizeOptions{
		Height: window.Height(),
		Width:  window.Width(),
	})
	if err != nil {
		return fmt.Errorf("failed to resize: %w", err)
	}

	return nil
}

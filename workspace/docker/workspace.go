package docker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

var (
	imageRef   string = os.Getenv("IMAGE_URL")
	imageUser  string = os.Getenv("IMAGE_USER")
	imageCmd   string = os.Getenv("IMAGE_CMD")
	createOpts        = types.ExecConfig{
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
	stopTimeout = 10 * time.Second
)

func containerName(userName values.UserName) string {
	return fmt.Sprintf("user-%s", userName)
}

type Workspace struct{}

func NewWorkspace() (*Workspace, error) {
	return &Workspace{}, nil
}

func (w *Workspace) Create(ctx context.Context, userName values.UserName) (*domain.Workspace, error) {
	ctnName := containerName(userName)
	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		User:  imageUser,
		Tty:   true,
	}, nil, nil, nil, ctnName)
	if errdefs.IsConflict(err) {
		ctnInfo, err := cli.ContainerInspect(ctx, ctnName)
		if err != nil {
			return nil, fmt.Errorf("failed to inspect container: %w", err)
		}

		workspaceID := values.NewWorkspaceID(ctnInfo.ID)
		workspaceName := values.NewWorkspaceName(ctnName)
		return domain.NewWorkspace(workspaceID, workspaceName, userName), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	workspaceID := values.NewWorkspaceID(res.ID)
	workspaceName := values.NewWorkspaceName(ctnName)

	return domain.NewWorkspace(workspaceID, workspaceName, userName), nil
}

func (w *Workspace) Start(ctx context.Context, workspace *domain.Workspace) error {
	err := cli.ContainerStart(ctx, string(workspace.ID()), types.ContainerStartOptions{})
	if err != nil && !errdefs.IsConflict(err) {
		return fmt.Errorf("failed to start container: %w", err)
	}
	workspace.Status = values.StatusUp

	return nil
}

func (w *Workspace) Stop(ctx context.Context, workspace *domain.Workspace) error {
	err := cli.ContainerStop(ctx, string(workspace.ID()), &stopTimeout)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	workspace.Status = values.StatusDown

	return nil
}

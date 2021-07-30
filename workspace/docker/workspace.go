package docker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	upLabel   = "up"
	downLabel = "down"
)

var (
	stopTimeout = 10 * time.Second
	cpuLimit    int64
	memoryLimit int64
)

var containerCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Help:      "Number of containers.",
	Namespace: "webshell",
	Name:      "container_num",
}, []string{"status"})

func containerName(userName values.UserName) string {
	return fmt.Sprintf("user-%s", userName)
}

type Workspace struct{}

func NewWorkspace() (*Workspace, error) {
	floatCPULimit, err := strconv.ParseFloat(os.Getenv("CPU_LIMIT"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cpu limit: %w", err)
	}
	cpuLimit = int64(floatCPULimit * 1e9)

	floatMemoryLimit, err := strconv.ParseFloat(os.Getenv("MEMORY_LIMIT"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid memory limit: %w", err)
	}
	memoryLimit = int64(floatMemoryLimit * 1e6)

	return &Workspace{}, nil
}

func (w *Workspace) Create(ctx context.Context, userName values.UserName) (*domain.Workspace, error) {
	ctnName := containerName(userName)
	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		User:  imageUser,
		Tty:   true,
	}, &container.HostConfig{
		Resources: container.Resources{
			NanoCPUs: cpuLimit,
			Memory:   memoryLimit,
		},
	}, nil, nil, ctnName)
	if errdefs.IsConflict(err) {
		ctnInfo, err := cli.ContainerInspect(ctx, ctnName)
		if err != nil {
			return nil, fmt.Errorf("failed to inspect container: %w", err)
		}

		workspaceID := values.NewWorkspaceID(ctnInfo.ID)
		workspaceName := values.NewWorkspaceName(ctnName)
		containerCounter.WithLabelValues(downLabel).Inc()

		return domain.NewWorkspace(workspaceID, workspaceName, userName), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	workspaceID := values.NewWorkspaceID(res.ID)
	workspaceName := values.NewWorkspaceName(ctnName)
	containerCounter.WithLabelValues(downLabel).Inc()

	return domain.NewWorkspace(workspaceID, workspaceName, userName), nil
}

func (w *Workspace) Start(ctx context.Context, workspace *domain.Workspace) error {
	err := cli.ContainerStart(ctx, string(workspace.ID()), types.ContainerStartOptions{})
	if err != nil && !errdefs.IsConflict(err) {
		return fmt.Errorf("failed to start container: %w", err)
	}
	workspace.Status = values.StatusUp
	containerCounter.WithLabelValues(downLabel).Dec()
	containerCounter.WithLabelValues(upLabel).Inc()

	return nil
}

func (w *Workspace) Stop(ctx context.Context, workspace *domain.Workspace) error {
	err := cli.ContainerStop(ctx, string(workspace.ID()), &stopTimeout)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	workspace.Status = values.StatusDown
	containerCounter.WithLabelValues(upLabel).Dec()
	containerCounter.WithLabelValues(downLabel).Inc()

	return nil
}

func (w *Workspace) Recreate(ctx context.Context, workspace *domain.Workspace) (*domain.Workspace, error) {
	err := cli.ContainerRemove(ctx, string(workspace.ID()), types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove container: %w", err)
	}
	containerCounter.WithLabelValues(upLabel).Dec()
	containerCounter.WithLabelValues(downLabel).Inc()

	userName := workspace.UserName()
	ctnName := string(workspace.Name())
	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		User:  imageUser,
		Tty:   true,
	}, &container.HostConfig{
		Resources: container.Resources{
			NanoCPUs: cpuLimit,
			Memory:   memoryLimit,
		},
	}, nil, nil, ctnName)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	containerCounter.WithLabelValues(downLabel).Dec()
	containerCounter.WithLabelValues(upLabel).Inc()

	workspaceID := values.NewWorkspaceID(res.ID)
	workspaceName := values.NewWorkspaceName(ctnName)

	return domain.NewWorkspace(workspaceID, workspaceName, userName), nil
}

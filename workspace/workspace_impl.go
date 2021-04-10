package workspace

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"golang.org/x/sync/errgroup"
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
	containerMap = sync.Map{}
	stopTimeout  = 10 * time.Second
)

type containerInfo struct {
	id         string
	manageChan chan struct{}
}

func containerName(userName values.UserName) string {
	return fmt.Sprintf("user-%s", userName)
}

type Workspace struct {
	cli *client.Client
}

func NewWorkspace() (*Workspace, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull image: %w", err)
	}
	io.Copy(os.Stdout, reader)

	return &Workspace{
		cli: cli,
	}, nil
}

func (w *Workspace) Create(ctx context.Context, userName values.UserName) error {
	ctnName := containerName(userName)
	res, err := w.cli.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		User:  imageUser,
	}, nil, nil, nil, ctnName)
	if errdefs.IsConflict(err) {
		ctnInfo, err := w.cli.ContainerInspect(ctx, ctnName)
		if err != nil {
			return fmt.Errorf("failed to inspect container: %w", err)
		}

		containerMap.Store(userName, &containerInfo{
			id:         ctnInfo.ID,
			manageChan: make(chan struct{}, 20),
		})

		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	containerMap.Store(userName, &containerInfo{
		id:         res.ID,
		manageChan: make(chan struct{}, 20),
	})

	return nil
}

func (w *Workspace) Connect(ctx context.Context, userName values.UserName, connection *domain.Connection) error {
	iContainerInfo, ok := containerMap.Load(userName)
	if !ok {
		return errors.New("load container info error")
	}
	ctnInfo, ok := iContainerInfo.(*containerInfo)
	if !ok {
		return errors.New("parse container info error")
	}

	if len(ctnInfo.manageChan) >= 20 {
		return errors.New("too many shell")
	}

	err := w.cli.ContainerStart(ctx, ctnInfo.id, types.ContainerStartOptions{})
	if err != nil && !errdefs.IsConflict(err) {
		return fmt.Errorf("failed to start container: %w", err)
	}
	ctnInfo.manageChan <- struct{}{}
	defer func(ctnInfo *containerInfo) {
		<-ctnInfo.manageChan
		if len(ctnInfo.manageChan) == 0 {
			ctx := context.Background()
			err := w.cli.ContainerStop(ctx, ctnInfo.id, &stopTimeout)
			if err != nil {
				log.Fatalf("failed to stop container:%+v", err)
			}
		}
	}(ctnInfo)

	idRes, err := w.cli.ContainerExecCreate(ctx, ctnInfo.id, createOpts)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if connection.IsTty() {
		go func(ctx context.Context) {
		WINDOW_RESIZE:
			for {
				select {
				case win := <-connection.WindowReceiver():
					err := w.cli.ContainerExecResize(ctx, idRes.ID, types.ResizeOptions{
						Height: win.Height(),
						Width:  win.Width(),
					})
					if err != nil {
						log.Println(err)
					}
				case <-ctx.Done():
					break WINDOW_RESIZE
				}
			}
		}(ctx)
	}

	stream, err := w.cli.ContainerExecAttach(ctx, idRes.ID, attachOpts)
	if err != nil {
		return fmt.Errorf("failed to attach container: %w", err)
	}
	defer stream.Close()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if connection.IsTty() {
			_, err = io.Copy(connection.Stdout(), stream.Reader)
			if err != nil {
				return fmt.Errorf("failed to copy stdout: %w", err)
			}
		} else {
			_, err = stdcopy.StdCopy(connection.Stdout(), connection.Stderr(), stream.Reader)
			if err != nil {
				return fmt.Errorf("failed to stdcopy: %w", err)
			}
		}

		return nil
	})

	eg.Go(func() error {
		defer stream.CloseWrite()
		_, err := io.Copy(stream.Conn, connection.Stdin())
		if err != nil {
			return fmt.Errorf("failed to copy stdin: %w", err)
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("failed to stdout: %w", err)
	}

	return nil
}

func (*Workspace) Remove(ctx context.Context, userName values.UserName) error {
	return nil
}

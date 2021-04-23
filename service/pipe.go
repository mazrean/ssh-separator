package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/store"
	"github.com/mazrean/separated-webshell/workspace"
)

var (
	welcome = os.Getenv("WELCOME")
)

type IPipe interface {
	Pipe(ctx context.Context, userName values.UserName, connection *domain.Connection) error
}

type Pipe struct {
	sw  store.IWorkspace
	wwc workspace.IWorkspaceConnection
	ww  workspace.IWorkspace
}

func NewPipe(sw store.IWorkspace, wwc workspace.IWorkspaceConnection, ww workspace.IWorkspace) *Pipe {
	return &Pipe{
		sw:  sw,
		wwc: wwc,
		ww:  ww,
	}
}

func (p *Pipe) Pipe(ctx context.Context, userName values.UserName, connection *domain.Connection) error {
	workspace, err := p.sw.Get(ctx, userName)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	if workspace.Status == values.StatusDown {
		err = p.ww.Start(ctx, workspace)
		if err != nil {
			return fmt.Errorf("failed to start workspace: %w", err)
		}
	}

	err = workspace.AddConnection()
	if err != nil {
		return fmt.Errorf("failed to add connection: %w", err)
	}

	workspaceConnection, err := p.wwc.Connect(ctx, workspace)
	if err != nil {
		return fmt.Errorf("connect to workspace error: %w", err)
	}
	defer func() {
		ctx := context.Background()

		err := p.wwc.Disconnect(ctx, workspaceConnection)
		if err != nil {
			log.Printf("failed to disconnect: %+v", err)
			return
		}

		err = workspace.RemoveConnection()
		if err != nil {
			log.Printf("connection num missmatch: %+v", err)
		}

		if workspace.ConnectionNum() == 0 {
			err = p.ww.Stop(ctx, workspace)
			if err != nil {
				log.Printf("failed to stop workspace: %+v", err)
			}
		}
	}()

	go func() {
		for win := range connection.WindowReceiver() {
			err := p.wwc.Resize(ctx, workspaceConnection, win)
			if err != nil {
				log.Printf("failed to resize window: %+v", err)
			}
		}
	}()

	go func() {
		defer connection.Close()
		if connection.IsTty() {
			if len(welcome) != 0 {
				_, err := io.Copy(connection.Stdout(), strings.NewReader(welcome))
				if err != nil {
					log.Printf("failed to copy fonts: %+v", err)
				}
			}

			_, err := io.Copy(connection.Stdout(), workspaceConnection.ReadCloser())
			if err != nil {
				log.Printf("failed to copy stdin: %+v\n", err)
			}
		} else {
			_, err := stdcopy.StdCopy(connection.Stdout(), connection.Stderr(), workspaceConnection.ReadCloser())
			if err != nil {
				log.Printf("failed to copy stdout: %+v\n", err)
			}
		}
	}()

	_, err = io.Copy(workspaceConnection.WriteCloser(), connection.Stdin())
	if err != nil {
		return fmt.Errorf("failed to copy stdin: %+v", err)
	}

	return nil
}

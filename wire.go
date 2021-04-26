//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mazrean/separated-webshell/api"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/repository/badger"
	"github.com/mazrean/separated-webshell/service"
	"github.com/mazrean/separated-webshell/ssh"
	"github.com/mazrean/separated-webshell/store"
	"github.com/mazrean/separated-webshell/store/gomap"
	"github.com/mazrean/separated-webshell/workspace"
	"github.com/mazrean/separated-webshell/workspace/docker"
)

var (
	transactionBind         = wire.Bind(new(repository.ITransaction), new(*badger.Transaction))
	storeWorkspaceBind      = wire.Bind(new(store.IWorkspace), new(*gomap.Workspace))
	repositoryUserBind      = wire.Bind(new(repository.IUser), new(*badger.User))
	workspaceBind           = wire.Bind(new(workspace.IWorkspace), new(*docker.Workspace))
	workspaceConnectionBind = wire.Bind(new(workspace.IWorkspaceConnection), new(*docker.WorkspaceConnection))
	serviceUserBind         = wire.Bind(new(service.IUser), new(*service.User))
	servicePipeBind         = wire.Bind(new(service.IPipe), new(*service.Pipe))
)

type Server struct {
	*service.Setup
	*api.API
	*ssh.SSH
}

func NewServer(setup *service.Setup, a *api.API, s *ssh.SSH) (*Server, error) {
	return &Server{
		Setup: setup,
		API:   a,
		SSH:   s,
	}, nil
}

func InjectServer() (*Server, func(), error) {
	wire.Build(
		NewServer,
		api.NewAPI,
		api.NewUser,
		gomap.NewWorkspace,
		badger.NewDB,
		badger.NewTransaction,
		badger.NewUser,
		service.NewSetup,
		service.NewUser,
		service.NewPipe,
		ssh.NewSSH,
		docker.NewWorkspace,
		docker.NewWorkspaceConnection,
		transactionBind,
		storeWorkspaceBind,
		repositoryUserBind,
		workspaceBind,
		workspaceConnectionBind,
		serviceUserBind,
		servicePipeBind,
	)

	return nil, nil, nil
}

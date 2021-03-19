//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mazrean/separated-webshell/api"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/service"
	"github.com/mazrean/separated-webshell/ssh"
	"github.com/mazrean/separated-webshell/workspace"
)

var (
	transactionBind    = wire.Bind(new(repository.ITransaction), new(*repository.Transaction))
	repositoryUserBind = wire.Bind(new(repository.IUser), new(*repository.User))
	workspaceBind      = wire.Bind(new(workspace.IWorkspace), new(*workspace.Workspace))
)

func InjectAPIServer() (*api.API, error) {
	wire.Build(
		api.NewAPI,
		api.NewUser,
		repository.NewTransaction,
		repository.NewUser,
		service.NewUser,
		//ssh.NewSSH,
		workspace.NewWorkspace,
		transactionBind,
		repositoryUserBind,
		workspaceBind,
	)

	return nil, nil
}

func InjectSSHServer() (*ssh.SSH, error) {
	wire.Build(
		//api.NewAPI,
		//api.NewUser,
		repository.NewTransaction,
		repository.NewUser,
		service.NewUser,
		ssh.NewSSH,
		workspace.NewWorkspace,
		transactionBind,
		repositoryUserBind,
		workspaceBind,
	)

	return nil, nil
}

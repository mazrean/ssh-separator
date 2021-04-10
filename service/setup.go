package service

import (
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/workspace"
)

type Setup struct {
	workspace.IWorkspace
	repository.IUser
}

func NewSetup(w workspace.IWorkspace, u repository.IUser) *Setup {
	return &Setup{
		IWorkspace: w,
		IUser:      u,
	}
}

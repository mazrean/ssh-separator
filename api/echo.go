package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	*User
}

func NewAPI(user *User) *API {
	return &API{
		User: user,
	}
}

func (api *API) Start(port int) error {
	e := echo.New()

	e.Use(middleware.Recover())

	e.POST("/new", api.User.PostNewUser)

	return e.Start(fmt.Sprintf(":%d", port))
}

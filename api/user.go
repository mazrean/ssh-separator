package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/service"
)

type User struct {
	*service.User
	*validator.Validate
}

func NewUser() *User {
	return &User{
		Validate: validator.New(),
	}
}

type postNewUserRequest struct {
	Name     string `json:"name" validate:"required,gt=0,lt=32,alphanum"`
	Password string `json:"password" validate:"required,gt=8,lt=32"`
}

func (u *User) PostNewUser(c echo.Context) error {
	req := postNewUserRequest{}
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	err = u.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = u.User.New(c.Request().Context(), &domain.User{
		Name:     req.Name,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserExist) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("user already exist"))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
	}

	return c.NoContent(http.StatusCreated)
}

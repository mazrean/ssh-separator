package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/service"
)

var (
	apiKey = os.Getenv("API_KEY")
)

type User struct {
	*service.User
	*validator.Validate
}

func NewUser(u *service.User) *User {
	return &User{
		User:     u,
		Validate: NewValidator(),
	}
}

type postNewUserRequest struct {
	APIKey   string `json:"key" validate:"required"`
	Name     string `json:"name" validate:"required,gt=0,lt=32,userName"`
	Password string `json:"cred" validate:"required,gt=8,lt=32,password"`
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

	if req.APIKey != apiKey {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
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

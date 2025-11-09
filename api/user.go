package api

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/service"
)

type Key string

type User struct {
	*service.User
	*validator.Validate
	apiKey string
}

func NewUser(u *service.User, apiKey Key) *User {
	return &User{
		User:     u,
		Validate: NewValidator(),
		apiKey:   string(apiKey),
	}
}

type postNewUserRequest struct {
	APIKey   string `json:"key" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"cred" validate:"required"`
}

func (u *User) PostNewUser(c echo.Context) error {
	req := postNewUserRequest{}
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	err = u.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if subtle.ConstantTimeCompare([]byte(req.APIKey), []byte(u.apiKey)) != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid api key")
	}

	userName, err := values.NewUserName(req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	password, err := values.NewPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = u.New(c.Request().Context(), userName, password)
	if errors.Is(err, service.ErrUserExist) {
		return echo.NewHTTPError(http.StatusBadRequest, "user already exist")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
	}

	return c.NoContent(http.StatusCreated)
}

type putResetRequest struct {
	APIKey string `json:"key" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

func (u *User) PutReset(c echo.Context) error {
	req := putResetRequest{}
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
	}

	err = u.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// セキュアなデフォルト: APIキーが空の場合は全てのリクエストを拒否
	if u.apiKey == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "API authentication is not properly configured")
	}

	if subtle.ConstantTimeCompare([]byte(req.APIKey), []byte(u.apiKey)) != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid api key")
	}

	userName, err := values.NewUserName(req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = u.ResetContainer(c.Request().Context(), userName)
	if errors.Is(err, service.ErrInvalidUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "no user")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to reset container: %w", err))
	}

	return nil
}

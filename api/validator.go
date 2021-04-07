package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	userNameExpression = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9_-]{0,14}[a-zA-Z0-9])?$`)
)

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("userName", userNameValidator)

	return v
}

func userNameValidator(fl validator.FieldLevel) bool {
	return userNameExpression.MatchString(fl.FieldName())
}

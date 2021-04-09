package values

import (
	"errors"
	"regexp"
)

type (
	UserName string
	Password string
	// HashedPassword 不明時""
	HashedPassword string
)

var (
	userNameExpression = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9_-]{0,14}[a-zA-Z0-9])?$`)
	passwordExpression = regexp.MustCompile(`^[a-zA-Z0-9]{8,32}$`)
)

func NewUserName(userName string) (UserName, error) {
	if !userNameExpression.MatchString(userName) {
		return "", errors.New("invalid user name")
	}

	return UserName(userName), nil
}

func NewPassword(password string) (Password, error) {
	if !passwordExpression.MatchString(password) {
		return "", errors.New("invalid password")
	}

	return Password(password), nil
}

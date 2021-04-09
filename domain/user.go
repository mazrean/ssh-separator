package domain

import "github.com/mazrean/separated-webshell/domain/values"

type User struct {
	name values.UserName
	values.Password
	values.HashedPassword
}

func NewUserWithPassword(name values.UserName, password values.Password) *User {
	return &User{
		name:     name,
		Password: password,
	}
}

func NewUserWithHashedPassword(name values.UserName, hashedPassword values.HashedPassword) *User {
	return &User{
		name:           name,
		HashedPassword: hashedPassword,
	}
}

func (u *User) GetName() values.UserName {
	return u.name
}

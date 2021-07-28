package domain

import "github.com/mazrean/separated-webshell/domain/values"

type User struct {
	name values.UserName
	values.HashedPassword
}

func NewUser(name values.UserName, hashedPassword values.HashedPassword) *User {
	return &User{
		name:           name,
		HashedPassword: hashedPassword,
	}
}

func (u *User) GetName() values.UserName {
	return u.name
}

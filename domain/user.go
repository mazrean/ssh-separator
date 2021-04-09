package domain

type (
	UserName string
	Password string
	// HashedPassword 不明時""
	HashedPassword string
)

type User struct {
	name UserName
	Password
	HashedPassword
}

func NewUserWithPassword(name UserName, password Password) *User {
	return &User{
		name:     name,
		Password: password,
	}
}

func NewUserWithHashedPassword(name UserName, hashedPassword HashedPassword) *User {
	return &User{
		name:           name,
		HashedPassword: hashedPassword,
	}
}

func (u *User) GetName() UserName {
	return u.name
}

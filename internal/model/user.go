package model

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
}

func NewUser(id uuid.UUID, email string) *User {
	return &User{Id: id, Email: email}
}

func (u *User) ChangeEmail(email string) {
	u.Email = email
}

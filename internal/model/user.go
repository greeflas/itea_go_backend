package model

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id           uuid.UUID `bun:",pk"`
	Email        string
	PasswordHash sql.NullString
}

func NewUser(id uuid.UUID, email string) *User {
	return &User{Id: id, Email: email}
}

func (u *User) ChangeEmail(email string) {
	u.Email = email
}

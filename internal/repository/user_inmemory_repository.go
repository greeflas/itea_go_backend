package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/greeflas/itea_go_backend/internal/model"
)

type UserInMemoryRepository struct {
	users []*model.User
}

func NewUserInMemoryRepository() *UserInMemoryRepository {
	return &UserInMemoryRepository{
		users: make([]*model.User, 0),
	}
}

func (r *UserInMemoryRepository) GetAll() []*model.User {
	return r.users
}

func (r *UserInMemoryRepository) Create(user *model.User) {
	r.users = append(r.users, user)
}

func (r *UserInMemoryRepository) GetById(id uuid.UUID) (*model.User, error) {
	for _, user := range r.users {
		if user.Id == id {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *UserInMemoryRepository) Delete(user *model.User) error {
	userIdx := -1

	for index, u := range r.users {
		if u.Id == user.Id {
			userIdx = index
			break
		}
	}

	if userIdx == -1 {
		return errors.New("user not found")
	}

	r.users = append(r.users[:userIdx], r.users[userIdx+1:]...)

	return nil
}

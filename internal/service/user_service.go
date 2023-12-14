package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/greeflas/itea_go_backend/internal/model"
	"github.com/greeflas/itea_go_backend/internal/repository"
)

type NewUserArgs struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type UpdatedUserArgs struct {
	Email string `json:"email"`
}

type UserService struct {
	userRepository *repository.UserBunRepository
}

func NewUserService(userRepository *repository.UserBunRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Create(ctx context.Context, args *NewUserArgs) error {
	id, err := uuid.Parse(args.Id)
	if err != nil {
		return err
	}

	user := model.NewUser(
		id,
		args.Email,
	)

	return s.userRepository.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, args *UpdatedUserArgs) error {
	user, err := s.userRepository.GetById(ctx, id)
	if err != nil {
		return err
	}

	user.ChangeEmail(args.Email)

	return s.userRepository.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepository.GetById(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepository.Delete(ctx, user)
}

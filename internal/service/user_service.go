package service

import (
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
	userRepository *repository.UserInMemoryRepository
}

func NewUserService(userRepository *repository.UserInMemoryRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Create(args *NewUserArgs) error {
	id, err := uuid.Parse(args.Id)
	if err != nil {
		return err
	}

	user := model.NewUser(
		id,
		args.Email,
	)

	s.userRepository.Create(user)

	return nil
}

func (s *UserService) Update(id uuid.UUID, args *UpdatedUserArgs) error {
	user, err := s.userRepository.GetById(id)
	if err != nil {
		return err
	}

	user.ChangeEmail(args.Email)

	return nil
}

func (s *UserService) Delete(id uuid.UUID) error {
	user, err := s.userRepository.GetById(id)
	if err != nil {
		return err
	}

	return s.userRepository.Delete(user)
}

package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"mova-backend/internal/database"
	"mova-backend/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, clerkID, email string) (database.User, error) {
	user, err := s.userRepo.Create(ctx, clerkID, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return database.User{}, ErrUserAlreadyExists
		}
		return database.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByClerkID(ctx context.Context, clerkID string) (database.User, error) {
	user, err := s.userRepo.GetByClerkID(ctx, clerkID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return database.User{}, ErrUserNotFound
		}
		return database.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return database.User{}, ErrUserNotFound
		}
		return database.User{}, err
	}
	return user, nil
}

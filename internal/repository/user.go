package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"mova-backend/internal/database"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository struct {
	queries *database.Queries
}

func NewUserRepository(queries *database.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, clerkID, email string) (database.User, error) {
	user, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		ClerkID: clerkID,
		Email:   email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.User{}, ErrUserAlreadyExists
		}
		return database.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByClerkID(ctx context.Context, clerkID string) (database.User, error) {
	user, err := r.queries.GetUserByClerkID(ctx, clerkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.User{}, ErrUserNotFound
		}
		return database.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.User{}, ErrUserNotFound
		}
		return database.User{}, err
	}
	return user, nil
}

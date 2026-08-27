package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(passwordHash string, password string) error
}

type RegisterUser struct {
	users  UserRepository
	hasher PasswordHasher
}

func NewRegisterUser(
	users UserRepository,
	hasher PasswordHasher,
) *RegisterUser {
	return &RegisterUser{
		users:  users,
		hasher: hasher,
	}
}

type RegisterInput struct {
	Email    string
	Password string
}

func (uc *RegisterUser) Execute(
	ctx context.Context,
	input RegisterInput,
) (*domain.User, error) {
	if err := validateRegisterInput(input); err != nil {
		return nil, err
	}

	_, err := uc.users.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Status:       domain.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func validateRegisterInput(input RegisterInput) error {
	if input.Email == "" {
		return domain.ErrInvalidEmail
	}

	if input.Password == "" {
		return domain.ErrInvalidPassword
	}

	return nil
}

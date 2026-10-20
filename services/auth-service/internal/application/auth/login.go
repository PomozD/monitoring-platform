package auth

import (
	"context"
	"errors"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/domain"
)

type LoginUser struct {
	users  UserRepository
	hasher PasswordHasher
}

func NewLoginUser(
	users UserRepository,
	hasher PasswordHasher,
) *LoginUser {
	return &LoginUser{
		users:  users,
		hasher: hasher,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

func (uc *LoginUser) Execute(ctx context.Context, input LoginInput) (*domain.User, error) {
	if input.Email == "" || input.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := uc.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	if user.Status != domain.UserStatusActive {
		return nil, domain.ErrInvalidCredentials
	}

	if err := uc.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

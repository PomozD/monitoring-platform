package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

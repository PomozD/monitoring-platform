package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PomozD/monitoring-platform/services/auth-service/internal/application/auth"
	"github.com/PomozD/monitoring-platform/services/auth-service/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

var _ auth.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	const query = `
		INSERT INTO users (id, email, password_hash, status)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Status,
	).Scan(
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "users_email_unique" {
			return domain.ErrUserAlreadyExists
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get user by id: %w", domain.ErrUserNotFound)
		}

		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			status,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get user by email: %w", domain.ErrUserNotFound)
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

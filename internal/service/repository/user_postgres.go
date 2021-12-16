package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const userTable = "users"

type userRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *userRepo {
	return &userRepo{pg}
}

// Create creates new user with email and password
func (r *userRepo) Create(ctx context.Context, dto entity.UserSensitiveData) (entity.User, error) {
	sql, args, err := r.Builder.
		Insert(userTable).
		Columns("email", "username", "password").
		Values(dto.Email(), dto.Username(), dto.PasswordHash()).
		Suffix("RETURNING id, created_at, is_active, is_archive").
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("r.Builder.Insert: %w", err)
	}

	u := entity.User{
		Email:    dto.Email(),
		Username: dto.Username(),
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			// SQL err handling by code
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", entity.ErrUserUniqueViolation)
			}

			// Return more detailed error message
			return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", errors.New(pgErr.Detail))
		}

		return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	u.UpdatedAt = u.CreatedAt

	return u, nil
}

// Archive sets is_archive to isArchive for user with id
func (r *userRepo) Archive(ctx context.Context, uid string, archive bool) (entity.User, error) {
	updatedAt := time.Now()

	sql, args, err := r.Builder.
		Update(userTable).
		Set("is_archive", archive).
		Set("updated_at", updatedAt).
		Where(sq.Eq{"id": uid, "is_archive": !archive}).
		Suffix("RETURNING email, username, created_at, is_active").
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("r.Builder.Update: %w", err)
	}

	var u entity.User

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.Email,
		&u.Username,
		&u.CreatedAt,
		&u.IsActive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", entity.ErrUserNotFound)
		}

		return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	u.ID = uid
	u.IsArchive = archive
	u.UpdatedAt = updatedAt

	return u, nil
}

// GetByID returns user data by id
func (r *userRepo) GetByID(ctx context.Context, uid string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("email, username, created_at, updated_at, is_active, is_archive").
		From(userTable).
		Where(sq.Eq{"id": uid, "is_active": true, "is_archive": false}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("r.Builder.Select: %w", err)
	}

	var u entity.User

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.Email,
		&u.Username,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", entity.ErrUserNotFound)
		}

		return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	u.ID = uid

	return u, nil
}

// GetByEmail returns user data by email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, username, password, created_at, updated_at, is_active, is_archive").
		From(userTable).
		Where(sq.Eq{"email": email, "is_active": true, "is_archive": false}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("r.Builder.Select: %w", err)
	}

	var u entity.User

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", entity.ErrUserNotFound)
		}

		return entity.User{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	u.Email = email

	return u, nil
}

package repo

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	// Build SQL query string
	sql, args, err := r.Builder.
		Insert("users").
		Columns("email", "password").
		Values(u.Email, u.EncryptedPassword).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	// Execute query
	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.CreatedAt); err != nil {
		return err
	}

	return nil
}

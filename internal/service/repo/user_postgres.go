package repo

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Create(ctx context.Context, u domain.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("email", "password", "created_at").
		Values(u.Email, u.EncryptedPassword, u.CreatedAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "UserRepo - Create - r.Builder")
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "UserRepo - Create - r.Pool.Exec")
	}

	return nil
}


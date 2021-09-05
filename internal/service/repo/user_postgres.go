package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const table = "users"

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Create creates new user with email and password
func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	sql, args, err := r.Builder.
		Insert(table).
		Columns("email", "password").
		Values(u.Email, u.EncryptedPassword).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.CreatedAt); err != nil {
		return err
	}

	return nil
}

// GetPassword returns user password by id
func (r *UserRepo) GetPassword(ctx context.Context, id int) (string, error) {
	sql, args, err := r.Builder.
		Select("password").
		From(table).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return "", err
	}

	var pwd string

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&pwd); err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("user with given id not found")
		}

		return "", err
	}

	return pwd, nil
}

func (r *UserRepo) UpdateState(ctx context.Context, u *domain.User) error {
	sql, args, err := r.Builder.
		Update(table).
		Set("is_active", u.IsActive).
		Where(sq.Eq{"id": u.ID, "is_active": !u.IsActive}).
		Suffix("RETURNING is_active").
		ToSql()
	if err != nil {
		return err
	}

	var isActive bool

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&isActive); err != nil {
		if err == pgx.ErrNoRows {
			var userState string

			switch u.IsActive {
			case true:
				userState = "deactivated"
			case false:
				userState = "activated"
			}

			return errors.New(fmt.Sprintf("%s user with given id not found", userState))
		}

		return err
	}

	if u.IsActive != isActive {
		return errors.New("user state did not change")
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	sql, args, err := r.Builder.Update(table).
		Set("first_name", u.FirstName).
		Set("last_name", u.LastName).
		Where(sq.Eq{"id": u.ID, "email": u.Email}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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

// Insert creates new user with email and password
func (r *UserRepo) Insert(ctx context.Context, u *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	sql, args, err := r.Builder.
		Insert(table).
		Columns("email", "password").
		Values(u.Email, u.EncryptedPassword).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	resp := domain.CreateUserResponse{Email: u.Email}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&resp.ID, &resp.CreatedAt); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			// SQL err handling by code
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, errors.New(fmt.Sprintf("user with email %s already exists", u.Email))
			}

			// Return more detailed error message
			return nil, errors.New(pgErr.Detail)
		}

		return nil, err
	}

	return &resp, nil
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
			return "", errors.New(fmt.Sprintf("user with id %d not found", id))
		}

		return "", err
	}

	return pwd, nil
}

func (r *UserRepo) UpdateState(ctx context.Context, resp *domain.UpdateStateUserResponse) error {
	sql, args, err := r.Builder.
		Update(table).
		Set("is_active", resp.IsActive).
		Set("updated_at", resp.UpdatedAt).
		Where(sq.Eq{"id": resp.ID, "is_active": !resp.IsActive}).
		Suffix("RETURNING is_active").
		ToSql()
	if err != nil {
		return err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&resp.IsActive); err != nil {
		// Create error message if activated/deactivated user not found
		if err == pgx.ErrNoRows {
			var userState string

			if resp.IsActive {
				userState = "deactivated"
			} else {
				userState = "activated"
			}

			return errors.New(fmt.Sprintf("%s user with id %d not found", userState, resp.ID))
		}

		return err
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	sql, args, err := r.Builder.Update(table).
		Set("first_name", u.FirstName).
		Set("last_name", u.LastName).
		Set("username", u.Username).
		Where(sq.Eq{"id": u.ID}).
		ToSql()
	if err != nil {
		return err
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			// SQL err handling by code
			if pgErr.Code == pgerrcode.UniqueViolation {
				return errors.New(fmt.Sprintf("user with username %s already exists", u.Username))
			}

			// Return more detailed error message
			return errors.New(pgErr.Detail)
		}

		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New(fmt.Sprintf("user with id %d not found", u.ID))
	}

	return nil
}

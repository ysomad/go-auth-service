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

func (r *UserRepo) Archive(ctx context.Context, resp *domain.ArchiveUserResponse) error {
	sql, args, err := r.Builder.
		Update(table).
		Set("is_archive", resp.IsArchive).
		Set("updated_at", resp.UpdatedAt).
		Where(sq.Eq{"id": resp.ID, "is_archive": !resp.IsArchive}).
		Suffix("RETURNING is_archive").
		ToSql()
	if err != nil {
		return err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&resp.IsArchive); err != nil {
		// Create error message if activated/deactivated user not found
		if err == pgx.ErrNoRows {
			var state string

			if resp.IsArchive {
				state = "archived"
			} else {
				state = "not archived"
			}

			return errors.New(fmt.Sprintf("%s user with id %d not found", state, resp.ID))
		}

		return err
	}

	return nil
}

func (r *UserRepo) PartialUpdate(ctx context.Context) error {

	return nil
}

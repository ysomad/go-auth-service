package repo

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const _userTable = "users"

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Create creates new user with email and password
func (r *UserRepo) Create(ctx context.Context, email string, password string) error {
	sql, args, err := r.Builder.
		Insert(_userTable).
		Columns("email", "password").
		Values(email, password).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			// SQL err handling by code
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entity.UserUniqueViolationErr
			}

			// Return more detailed error message
			return errors.New(pgErr.Detail)
		}

		return err
	}

	return nil
}

// Archive sets is_archive to isArchive for user with id
func (r *UserRepo) Archive(ctx context.Context, id int, isArchive bool) error {
	sql, args, err := r.Builder.
		Update(_userTable).
		Set("is_archive", isArchive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id, "is_archive": !isArchive}).
		ToSql()
	if err != nil {
		return err
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return entity.UserNotFoundErr
	}

	return nil
}

// PartialUpdate update User column values with values presented in cols
func (r *UserRepo) PartialUpdate(ctx context.Context, id int, cols map[string]interface{}) error {
	sql, args, err := r.Builder.
		Update(_userTable).
		SetMap(cols).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id, "is_archive": false}).
		ToSql()
	if err != nil {
		return err
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entity.UserUniqueViolationErr
			}

			return errors.New(pgErr.Detail)
		}

		return err
	}
	if ct.RowsAffected() == 0 {
		return entity.UserNotFoundErr
	}

	return nil
}

// GetByID returns user data by its id
func (r *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	u := entity.User{ID: id}

	sql, args, err := r.Builder.
		Select("email, username, first_name, last_name, created_at, updated_at, is_active, is_archive").
		From(_userTable).
		Where(sq.Eq{"id": u.ID, "is_active": true, "is_archive": false}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.Email,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.UserNotFoundErr
		}

		return nil, err
	}

	return &u, nil
}

// GetByEmail returns user data by its email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u := entity.User{Email: email}

	sql, args, err := r.Builder.
		Select("id, email, password, is_active, is_archive").
		From(_userTable).
		Where(sq.Eq{"email": u.Email, "is_active": true, "is_archive": false}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.UserIncorrectCredsErr
		}

		return nil, err
	}

	return &u, nil
}

package repo

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

const table = "users"

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Create creates new user with email and password
func (r *UserRepo) Create(ctx context.Context, email string, password string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Insert(table).
		Columns("email", "password").
		Values(email, password).
		Suffix("RETURNING id, created_at, updated_at, is_active, is_archive").
		ToSql()
	if err != nil {
		return nil, err
	}

	u := entity.User{Email: email}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
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

	return &u, nil
}

// Archive sets is_archive to isArchive for user with id
func (r *UserRepo) Archive(ctx context.Context, id int, isArchive bool) error {
	sql, args, err := r.Builder.
		Update(table).
		Set("is_archive", isArchive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id, "is_archive": !isArchive}).
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	// Create error message if archived/not archived user not found
	if commandTag.RowsAffected() == 0 {
		var state string

		if !isArchive {
			state = "archived"
		} else {
			state = "not archived"
		}

		return errors.New(fmt.Sprintf("%s user with id %d not found", state, id))
	}

	return nil
}

// PartialUpdate update User column values with values presented in cols
func (r *UserRepo) PartialUpdate(ctx context.Context, id int, cols map[string]interface{}) (*entity.User, error) {


	u := entity.User{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	sql, args, err := r.Builder.
		Update(table).
		SetMap(cols).
		Set("updated_at", u.UpdatedAt).
		Where(sq.Eq{"id": u.ID, "is_archive": false}).
		Suffix("RETURNING username, first_name, last_name, email, created_at, is_active, is_archive").
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.CreatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("user with id %d not found", u.ID))
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, errors.New(fmt.Sprintf("user with username %s already exists", *u.Username))
			}

			return nil, errors.New(pgErr.Detail)
		}

		return nil, err
	}

	return &u, nil
}

// GetByID returns user data by its id
func (r *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	u := entity.User{ID: id}

	sql, args, err := r.Builder.
		Select("email, username, first_name, last_name, created_at, updated_at, is_active, is_archive").
		From(table).
		Where(sq.Eq{"id": u.ID}).
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
			return nil, errors.New(fmt.Sprintf("user with id %d not found", u.ID))
		}

		return nil, err
	}

	return &u, nil
}

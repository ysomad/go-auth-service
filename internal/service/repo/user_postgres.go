package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const (
	table = "users"
)

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
		Suffix("RETURNING id, created_at, updated_at, is_active, is_archive").
		ToSql()
	if err != nil {
		return err
	}

	if err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.IsActive, &u.IsArchive); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			// SQL err handling by code
			if pgErr.Code == pgerrcode.UniqueViolation {
				return errors.New(fmt.Sprintf("user with email %s already exists", u.Email))
			}

			// Return more detailed error message
			return errors.New(pgErr.Detail)
		}

		return err
	}

	return nil
}

func (r *UserRepo) Archive(ctx context.Context, req *domain.ArchiveUserRequest) error {
	sql, args, err := r.Builder.
		Update(table).
		Set("is_archive", req.IsArchive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.ID, "is_archive": !*req.IsArchive}).
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

		if !*req.IsArchive {
			state = "archived"
		} else {
			state = "not archived"
		}

		return errors.New(fmt.Sprintf("%s user with id %d not found", state, req.ID))
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	userMap, err := stripNilValues(map[string]interface{}{
		"username":   *u.Username,
		"first_name": *u.FirstName,
		"last_name":  *u.LastName,
	})
	if err != nil {
		return err
	}

	sql, args, err := r.Builder.
		Update(table).
		SetMap(userMap).
		Set("updated_at", u.UpdatedAt).
		Where(sq.Eq{"id": u.ID}).
		Suffix("RETURNING username, first_name, last_name, email, created_at, is_active, is_archive").
		ToSql()
	if err != nil {
		return err
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		// Nullable values
		&u.Username,
		&u.FirstName,
		&u.LastName,

		// Not null values
		&u.Email,
		&u.CreatedAt,
		&u.IsActive,
		&u.IsArchive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return errors.New(fmt.Sprintf("user with id %d not found", u.ID))
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return errors.New(fmt.Sprintf("user with username %s already exists", *u.Username))
			}

			return errors.New(pgErr.Detail)
		}

		return err
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, u *domain.User) error {
	sql, args, err := r.Builder.
		Select("email, username, first_name, last_name, created_at, updated_at, is_active, is_archive").
		From(table).
		Where(sq.Eq{"id": u.ID}).
		ToSql()
	if err != nil {
		return err
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
			return errors.New(fmt.Sprintf("user with id %d not found", u.ID))
		}

		return err
	}

	return nil
}

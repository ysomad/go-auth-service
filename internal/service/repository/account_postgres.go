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

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const _accTable = "accounts"

type accountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) *accountRepo {
	return &accountRepo{pg}
}

func (r *accountRepo) Create(ctx context.Context, cred domain.AccountCredentials) error {
	// TODO: refactor
	// TODO: generic error pkg/httperror/

	sql, args, err := r.Builder.
		Insert(_accTable).
		Columns("email", "password").
		Values(cred.Email(), cred.PasswordHash()).
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
				return domain.ErrUserUniqueViolation
			}

			// Return more detailed error message
			return errors.New(pgErr.Detail)
		}

		return err
	}

	return nil
}

func (r *accountRepo) FindByID(ctx context.Context, aid string) (domain.Account, error) {
	// TODO: refactor
	// TODO: generic error pkg/httperror/

	sql, args, err := r.Builder.
		Select("email, password, created_at, updated_at").
		From(_accTable).
		Where(sq.Eq{"id": aid, "is_archive": false}).
		ToSql()
	if err != nil {
		return domain.Account{}, fmt.Errorf("r.Builder.Select: %w", err)
	}

	acc := domain.Account{ID: aid}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.Email,
		&acc.PasswordHash,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", domain.ErrUserNotFound)
		}

		return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return acc, nil
}

func (r *accountRepo) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
	// TODO: refactor
	// TODO: generic error pkg/httperror/

	sql, args, err := r.Builder.
		Select("id, password, created_at, updated_at").
		From(_accTable).
		Where(sq.Eq{"email": email, "is_archive": false}).
		ToSql()
	if err != nil {
		return domain.Account{}, fmt.Errorf("r.Builder.Select: %w", err)
	}

	acc := domain.Account{Email: email}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.ID,
		&acc.PasswordHash,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", domain.ErrUserNotFound)
		}

		return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return acc, nil
}

func (r *accountRepo) Archive(ctx context.Context, aid string, archive bool) error {
	// TODO: refactor
	// TODO: generic errors pkg/httperror

	sql, args, err := r.Builder.
		Update(_accTable).
		Set("is_archive", archive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": aid, "is_archive": !archive}).
		ToSql()
	if err != nil {
		return err
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

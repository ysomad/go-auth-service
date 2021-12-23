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

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const _accTable = "accounts"

type accountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) *accountRepo {
	return &accountRepo{pg}
}

func (r *accountRepo) Create(ctx context.Context, acc domain.Account) error {
	sql, args, err := r.Builder.
		Insert(_accTable).
		Columns("email", "password").
		Values(acc.Email, acc.PasswordHash).
		ToSql()
	if err != nil {
		return fmt.Errorf("r.Builder.Insert: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("r.Pool.Exec: %w", apperrors.ErrAccountAlreadyExist)
			}
		}

		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *accountRepo) FindByID(ctx context.Context, aid string) (domain.Account, error) {
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
			return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", apperrors.ErrAccountNotFound)
		}

		return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return acc, nil
}

func (r *accountRepo) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
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
			return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", apperrors.ErrAccountNotFound)
		}

		return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return acc, nil
}

func (r *accountRepo) Archive(ctx context.Context, aid string, archive bool) error {
	sql, args, err := r.Builder.
		Update(_accTable).
		Set("is_archive", archive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": aid, "is_archive": !archive}).
		ToSql()
	if err != nil {
		return fmt.Errorf("r.Builder.Update: %w", err)
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", apperrors.ErrAccountNotFound)
	}

	return nil
}

package repo

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

const defaultEntityCap = 64

type TranslationRepo struct {
	*postgres.Postgres
}

func NewTranslationRepo(pg *postgres.Postgres) *TranslationRepo {
	return &TranslationRepo{pg}
}

func (r *TranslationRepo) GetHistory(ctx context.Context) ([]domain.Translation, error) {
	sql, _, err := r.Builder.
		Select("source, destination, original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "TranslationRepo - GetHistory - r.Builder")
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "TranslationRepo - GetHistory - r.Pool.Query")
	}
	defer rows.Close()

	entities := make([]domain.Translation, 0, defaultEntityCap)

	for rows.Next() {
		e := domain.Translation{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, errors.Wrap(err, "TranslationRepo - GetHistory - rows.Scan")
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r *TranslationRepo) Store(ctx context.Context, entity domain.Translation) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(entity.Source, entity.Destination, entity.Original, entity.Translation).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "TranslationRepo - Store - r.Builder")
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "TranslationRepo - Store - r.Pool.Exec")
	}

	return nil
}

package models

import (
	"context"
	"database/sql"
	"fmt"
	"rashikzaman/api/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type QueryParam struct {
	Relations  []string
	Pagination utils.PaginationConfig
	Alias      string
}

func Create(ctx context.Context, db bun.IDB, model interface{}) error {
	_, err := db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}

func Select(ctx context.Context, db bun.IDB, models interface{}, queryParam QueryParam) (int, error) {
	query := db.NewSelect().
		Model(models)

	count, err := queryParam.Pagination.BuildPaginationQuery(ctx, query)
	if err != nil {
		return 0, errors.Wrap(err, err.Error())
	}

	if len(queryParam.Relations) != 0 {
		for _, relation := range queryParam.Relations {
			query.Relation(relation)
		}
	}

	err = query.Scan(ctx)
	if err != nil {
		return 0, errors.Wrap(err, err.Error())
	}

	return count, nil
}

func SelectAll(ctx context.Context, db bun.IDB, models interface{}, queryParam QueryParam) error {
	query := db.NewSelect().
		Model(models)

	if len(queryParam.Relations) != 0 {
		for _, relation := range queryParam.Relations {
			query.Relation(relation)
		}
	}

	err := query.Scan(ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}

func SelectByID(ctx context.Context, db bun.IDB, id uuid.UUID, model interface{}, queryParam QueryParam) error {
	query := db.NewSelect().
		Model(model)

	if queryParam.Alias != "" {
		query.Where(queryParam.Alias+".id = ?", id)
	} else {
		query.Where("id = ?", id)
	}

	if len(queryParam.Relations) != 0 {
		for _, relation := range queryParam.Relations {
			query.Relation(relation)
		}
	}

	err := query.Scan(ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}

func Update(ctx context.Context, db bun.IDB, model interface{}) error {
	query := db.NewUpdate().
		Model(model)

	_, err := query.WherePK().Exec(ctx)

	return err
}

func Delete(ctx context.Context, db bun.IDB, model interface{}) error {
	_, err := db.NewDelete().
		Model(model).
		WherePK().
		Exec(ctx)

	return err
}

func WithTransaction(ctx context.Context, db bun.IDB, fn func(*bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(&tx)
	if err != nil {
		txError := tx.Rollback()
		if txError != nil {
			return fmt.Errorf("transaction rollback error: %w", txError)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		txError := tx.Rollback()
		if txError != nil {
			return fmt.Errorf("transaction rollback error: %w", txError)
		}

		return err
	}

	return nil
}

func WithRollBackOnlyTransaction(ctx context.Context, db bun.IDB, fn func(*bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(&tx)
	if err != nil {
		txError := tx.Rollback()
		if txError != nil {
			return fmt.Errorf("transaction rollback error: %w", txError)
		}

		return err
	}

	txError := tx.Rollback()
	if txError != nil {
		return fmt.Errorf("transaction rollback error: %w", txError)
	}

	return nil
}

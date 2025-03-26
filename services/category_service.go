package services

import (
	"context"
	"rashikzaman/api/models"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func FetchCategories(ctx context.Context, db bun.IDB) ([]models.Category, error) {
	categories := []models.Category{}

	err := models.SelectAll(ctx, db, &categories, models.QueryParam{})
	if err != nil {
		return categories, errors.Wrap(err, err.Error())
	}

	return categories, nil
}

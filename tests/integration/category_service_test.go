package integration_test

import (
	"context"
	"rashikzaman/api/models"
	"rashikzaman/api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func (s *TestSuite) TestFetchCategories() {
	ctx := context.Background()

	// Create test data
	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		categories, err := services.FetchCategories(ctx, tx)
		require.NoError(s.T(), err)

		require.Equal(s.T(), 1, len(categories))
		require.Equal(s.T(), categories[0].ID, category.ID)

		return nil
	})

	require.NoError(s.T(), err)
}

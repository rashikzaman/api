package services

import (
	"context"
	"rashikzaman/api/models"
	"rashikzaman/api/utils"

	"github.com/uptrace/bun"
)

func FetchSkills(ctx context.Context, db bun.IDB) ([]string, error) {
	var uniqueSkills []string

	// Query to select distinct values from the RequiredSkills array
	err := db.NewSelect().
		ColumnExpr("DISTINCT unnest(required_skills) AS skill").
		Model((*models.Task)(nil)).
		Scan(ctx, &uniqueSkills)

	return utils.DeleteEmptyFromSlice(uniqueSkills), err
}

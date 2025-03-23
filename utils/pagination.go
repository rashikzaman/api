package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type PaginationConfig struct {
	Page            int
	Limit           int
	Offset          int
	SortColumn      string
	SortDirection   string
	Search          string
	ColumnSearch    map[string]string
	SelectedColumns []string
	DateField       string
	StartDate       string
	EndDate         string
}

func (paginationConfig *PaginationConfig) BuildPaginationQueryWithSorting(
	ctx context.Context, query *bun.SelectQuery, distinctCountColumn ...string,
) (int, error) {
	count, err := paginationConfig.BuildPaginationQuery(ctx, query, distinctCountColumn...)
	if err != nil {
		return 0, err
	}

	if paginationConfig.SortColumn == "" {
		paginationConfig.SortColumn = "created_at"
	}

	if paginationConfig.SortDirection == "" {
		paginationConfig.SortDirection = "DESC"
	}

	query.Order(paginationConfig.SortColumn + " " + paginationConfig.SortDirection)

	return count, nil
}

func (paginationConfig *PaginationConfig) BuildPaginationQuery(
	ctx context.Context, query *bun.SelectQuery, distinctCountColumn ...string,
) (count int, err error) {
	if len(distinctCountColumn) > 0 {
		countQuery := *query
		countQuery.ExcludeColumn("*").ColumnExpr(fmt.Sprintf("COUNT(DISTINCT %s)", distinctCountColumn[0]))
	}

	count, err = query.Count(ctx)
	if err != nil {
		return 0, err
	}

	if paginationConfig.Page > 0 && paginationConfig.Limit > 0 {
		query.Limit(paginationConfig.Limit)
		query.Offset((paginationConfig.Page - 1) * paginationConfig.Limit)
	} else {
		if paginationConfig.Limit > 0 {
			query.Limit(paginationConfig.Limit)
		}

		if paginationConfig.Offset > 0 {
			query.Offset(paginationConfig.Offset)
		}
	}

	return count, nil
}

func PaginationConfigFromRequest(c *gin.Context) PaginationConfig {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("per_page", "25"))
	if err != nil || limit < 1 {
		limit = 25
	}

	offset := (page - 1) * limit

	return PaginationConfig{
		Page:          page,
		Limit:         limit,
		Offset:        offset,
		SortColumn:    c.Query("sort-column"),
		SortDirection: c.Query("sort-direction"),
	}
}

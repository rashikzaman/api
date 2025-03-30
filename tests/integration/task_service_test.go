package integration_test

import (
	"context"
	"rashikzaman/api/models"
	"rashikzaman/api/services"
	"testing"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func (s *TestSuite) TestCreateTask() {
	ctx := context.Background()

	// Create test user

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {

		user := &models.User{Base: models.Base{ID: uuid.New()}}
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		// Create test category
		category := &models.Category{
			Base: models.Base{
				ID: uuid.New(),
			},
			Name: "Test Category",
		}
		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		// Test data
		taskBody := &models.Task{
			Title:                   "Test Task",
			Description:             "Test Description",
			RequiredVolunteersCount: 5,
			RequiredSkills:          []string{"skill1", "skill2"},
			Latitude:                40.7128,
			Longitude:               -74.0060,
			FormattedAddress:        "New York, NY",
			CategoryID:              category.ID,
			Media: []*models.TaskMedia{
				{
					Base64: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==",
				},
			},
		}

		// Test function
		err = services.CreateTask(ctx, tx, taskBody, user.ID)
		require.NoError(s.T(), err)

		// Verify task was created
		var task models.Task
		err = tx.NewSelect().Model(&task).Where("id = ?", taskBody.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "Test Task", task.Title)
		require.Equal(s.T(), user.ID, task.UserID)

		// Verify media was created
		var media []models.TaskMedia
		err = tx.NewSelect().Model(&media).Where("task_id = ?", taskBody.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.Len(s.T(), media, 1)

		return nil
	})

	require.NoError(s.T(), err)
}

func (s *TestSuite) TestFetchTasks() {
	ctx := context.Background()

	// Create test data
	user := &models.User{Base: models.Base{ID: uuid.New()}}
	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}

	tasks := []models.Task{
		{
			Title:       "Task 1",
			Description: "Description 1",
			Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
			UserID:      user.ID,
			CategoryID:  category.ID,
		},
		{
			Title:       "Task 2",
			Description: "Description 2",
			Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0061, 40.7129}, SRID: 4326},
			UserID:      user.ID,
			CategoryID:  category.ID,
		},
	}

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		// Insert test data
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(&tasks).Exec(ctx)
		require.NoError(s.T(), err)

		// Test cases
		tests := []struct {
			name     string
			filter   services.Filter
			expected int
		}{
			{
				name:     "No filters",
				filter:   services.Filter{},
				expected: 2,
			},
			{
				name: "Filter by category",
				filter: services.Filter{
					CategoryIDs: []string{category.ID.String()},
				},
				expected: 2,
			},
			{
				name: "Filter by location",
				filter: services.Filter{
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				expected: 2,
			},
			{
				name: "Filter by search term",
				filter: services.Filter{
					SearchTerm: "Task 1",
				},
				expected: 1,
			},
		}

		for _, tt := range tests {
			s.T().Run(tt.name, func(t *testing.T) {
				result, count, err := services.FetchTasks(ctx, tx, models.QueryParam{}, tt.filter)
				require.NoError(t, err)
				require.Equal(t, tt.expected, len(result))
				require.Equal(t, tt.expected, count)
			})
		}

		return nil
	})

	require.NoError(s.T(), err)
}

func (s *TestSuite) TestFetchTaskByID() {
	ctx := context.Background()

	user := &models.User{Base: models.Base{ID: uuid.New()}}

	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}

	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
		UserID:      user.ID,
		CategoryID:  category.ID,
	}

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(task).Exec(ctx)
		require.NoError(s.T(), err)

		// Test function
		fetchedTask, err := services.FetchTaskByID(ctx, tx, task.ID, models.QueryParam{})
		require.NoError(s.T(), err)
		require.Equal(s.T(), task.ID, fetchedTask.ID)
		require.Equal(s.T(), "Test Task", fetchedTask.Title)

		return nil
	})

	require.NoError(s.T(), err)
}

func (s *TestSuite) TestUpdateTask() {
	ctx := context.Background()

	user := &models.User{Base: models.Base{ID: uuid.New()}}

	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}

	originalTask := &models.Task{
		Title:       "Original Title",
		Description: "Original Description",
		Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
		UserID:      user.ID,
		CategoryID:  category.ID,
	}

	models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(originalTask).Exec(ctx)
		require.NoError(s.T(), err)

		// Update data
		updatedTask := models.Task{
			Title:       "Updated Title",
			Description: "Updated Description",
			Latitude:    40.7130,
			Longitude:   -74.0062,
			CategoryID:  originalTask.CategoryID,
			UserID:      user.ID,
		}

		// Test function
		result, err := services.UpdateTask(ctx, tx, *originalTask, updatedTask)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "Updated Title", result.Title)

		// Verify in database
		var dbTask models.Task
		err = tx.NewSelect().Model(&dbTask).Where("id = ?", originalTask.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "Updated Title", dbTask.Title)

		return nil
	})
}

func (s *TestSuite) TestDeleteTask() {
	ctx := context.Background()

	// Create test data
	user := &models.User{Base: models.Base{ID: uuid.New()}}
	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}

	task := &models.Task{
		Title:       "Task to delete",
		Description: "Will be deleted",
		Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
		UserID:      user.ID,
		CategoryID:  category.ID,
	}

	models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		// Insert test data
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(task).Exec(ctx)
		require.NoError(s.T(), err)

		// Test function
		err = services.DeleteTask(ctx, tx, task.ID)
		require.NoError(s.T(), err)

		// Verify deletion
		count, err := tx.NewSelect().Model((*models.Task)(nil)).Where("id = ?", task.ID).Count(ctx)
		require.NoError(s.T(), err)
		require.Equal(s.T(), 0, count)

		return nil
	})
}

func (s *TestSuite) TestApplyAndWithdrawFromTask() {
	ctx := context.Background()

	// Create test data
	user := &models.User{Base: models.Base{ID: uuid.New()}}
	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}
	task := &models.Task{
		Title:       "Task for application",
		Description: "Apply to this task",
		Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
		UserID:      user.ID,
		CategoryID:  category.ID,
	}

	models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(task).Exec(ctx)
		require.NoError(s.T(), err)

		// Test ApplyToTask
		err = services.ApplyToTask(ctx, tx, task.ID, user.ID)
		require.NoError(s.T(), err)

		// Verify application
		var userTask models.UserTask
		err = tx.NewSelect().Model(&userTask).
			Where("user_id = ?", user.ID).
			Where("task_id = ?", task.ID).
			Scan(ctx)
		require.NoError(s.T(), err)

		// Test WithdrawFromTask
		err = services.WithdrawFromTask(ctx, tx, task.ID, user.ID)
		require.NoError(s.T(), err)

		// Verify withdrawal
		count, err := tx.NewSelect().Model((*models.UserTask)(nil)).
			Where("user_id = ?", user.ID).
			Where("task_id = ?", task.ID).
			Count(ctx)
		require.NoError(s.T(), err)
		require.Equal(s.T(), 0, count)

		return nil
	})
}

func (s *TestSuite) TestFetchSubscribersForTask() {
	ctx := context.Background()

	// Create test data
	user1 := &models.User{Base: models.Base{ID: uuid.New()}}
	user2 := &models.User{Base: models.Base{ID: uuid.New()}}
	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}
	task := &models.Task{
		Base:        models.Base{ID: uuid.New()},
		Title:       "Task with subscribers",
		Description: "Multiple users applied",
		Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
		UserID:      user1.ID,
		CategoryID:  category.ID,
	}

	userTasks := []models.UserTask{
		{UserID: user1.ID, TaskID: task.ID, Base: models.Base{ID: uuid.New()}},
		{UserID: user2.ID, TaskID: task.ID, Base: models.Base{ID: uuid.New()}},
	}

	models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(user1).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(user2).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(task).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(&userTasks).Exec(ctx)
		require.NoError(s.T(), err)

		// Test function
		subscribers, err := services.FetchSubscribersForTask(ctx, tx, task.ID)
		require.NoError(s.T(), err)
		require.Len(s.T(), subscribers, 2)

		// Verify user IDs
		userIDs := make([]uuid.UUID, len(subscribers))
		for i, sub := range subscribers {
			userIDs[i] = sub.UserID
		}
		require.Contains(s.T(), userIDs, user1.ID)
		require.Contains(s.T(), userIDs, user2.ID)

		return nil
	})
}

func (s *TestSuite) TestApplyActionToTask() {
	ctx := context.Background()

	// Create test data
	user := &models.User{Base: models.Base{ID: uuid.New()}}
	category := &models.Category{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name: "Test Category",
	}
	task := &models.Task{
		Base:        models.Base{ID: uuid.New()},
		Title:       "Task to block",
		Description: "Will be blocked/unblocked",
		Location:    models.PostgisGeometry{Geometry: orb.Point{-74.0060, 40.7128}, SRID: 4326},
		Blocked:     false,
		UserID:      user.ID,
		CategoryID:  category.ID,
	}

	models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(category).Exec(ctx)
		require.NoError(s.T(), err)

		_, err = tx.NewInsert().Model(task).Exec(ctx)
		require.NoError(s.T(), err)

		// Test block action
		updatedTask, err := services.ApplyActionToTask(ctx, tx, task.ID, "block")
		require.NoError(s.T(), err)
		require.True(s.T(), updatedTask.Blocked)

		// Verify in database
		var dbTask models.Task
		err = tx.NewSelect().Model(&dbTask).Where("id = ?", task.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.True(s.T(), dbTask.Blocked)

		// Test unblock action
		updatedTask, err = services.ApplyActionToTask(ctx, tx, task.ID, "unblock")
		require.NoError(s.T(), err)
		require.False(s.T(), updatedTask.Blocked)

		// Verify in database
		err = tx.NewSelect().Model(&dbTask).Where("id = ?", task.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.False(s.T(), dbTask.Blocked)

		return nil
	})
}

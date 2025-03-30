package integration_test

import (
	"context"
	"database/sql"
	"errors"
	"rashikzaman/api/models"
	"rashikzaman/api/services"
	"rashikzaman/api/utils"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func (s *TestSuite) TestCreateUserFromClerk() {
	ctx := context.Background()

	testCases := []struct {
		name      string
		clerkID   string
		email     string
		firstName string
		lastName  string
		birthday  string
		wantError bool
	}{
		{
			name:      "success with all fields",
			clerkID:   "clerk_123",
			email:     "test@example.com",
			firstName: "John",
			lastName:  "Doe",
			birthday:  strconv.FormatInt(time.Now().UnixMilli(), 10),
			wantError: false,
		},
		{
			name:      "success without birthday",
			clerkID:   "clerk_456",
			email:     "test2@example.com",
			firstName: "Jane",
			lastName:  "Smith",
			birthday:  "",
			wantError: false,
		},
		{
			name:      "invalid birthday format",
			clerkID:   "clerk_789",
			email:     "test3@example.com",
			firstName: "Invalid",
			lastName:  "Birthday",
			birthday:  "not-a-number",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
				user, err := services.CreateUserFromClerk(ctx, tx, tc.clerkID, tc.email, tc.firstName, tc.lastName, tc.birthday)

				if tc.wantError {
					require.Error(t, err)
					return nil
				}

				require.NoError(t, err)
				require.Equal(t, tc.clerkID, user.ClerkID)
				require.Equal(t, tc.email, *user.Email)
				require.Equal(t, tc.firstName, user.FirstName)
				require.Equal(t, tc.lastName, user.LastName)
				require.Equal(t, "user", user.Role)

				if tc.birthday != "" {
					require.NotNil(t, user.DateOfBirth)
				} else {
					require.Nil(t, user.DateOfBirth)
				}

				// check if the user was actually created in the database
				var dbUser models.User
				err = tx.NewSelect().Model(&dbUser).Where("id = ?", user.ID).Scan(ctx)
				require.NoError(t, err)
				require.Equal(t, user.ID, dbUser.ID)

				return nil
			})
			require.NoError(t, err)
		})
	}
}

func (s *TestSuite) TestFetchUsersForAdmin() {
	ctx := context.Background()

	// Create test users
	users := []models.User{
		{
			ClerkID:   "clerk_1",
			Email:     ptr("user1@example.com"),
			FirstName: "User1",
			LastName:  "Test",
			Role:      "user",
		},
		{
			ClerkID:   "clerk_2",
			Email:     ptr("user2@example.com"),
			FirstName: "User2",
			LastName:  "Test",
			Role:      "user",
		},
		{
			ClerkID:   "clerk_3",
			Email:     ptr("user3@example.com"),
			FirstName: "User3",
			LastName:  "Test",
			Role:      "admin",
		},
	}

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		// Insert test data
		_, err := tx.NewInsert().Model(&users).Exec(ctx)
		require.NoError(s.T(), err)

		// Test cases
		tests := []struct {
			name       string
			queryParam models.QueryParam
			expected   int
		}{
			{
				name:       "fetch all users",
				queryParam: models.QueryParam{},
				expected:   3,
			},
			{
				name: "fetch with pagination",
				queryParam: models.QueryParam{
					Pagination: utils.PaginationConfig{
						Page:  1,
						Limit: 2,
					},
				},
				expected: 2,
			},
		}

		for _, tt := range tests {
			s.T().Run(tt.name, func(t *testing.T) {
				result, count, err := services.FetchUsersForAdmin(ctx, tx, tt.queryParam)
				require.NoError(t, err)
				require.Equal(t, tt.expected, len(result))
				require.Equal(t, 3, count) // Total count should always be 3
			})
		}

		return nil
	})
	require.NoError(s.T(), err)
}

func (s *TestSuite) TestGetUserByID() {
	ctx := context.Background()

	testUser := &models.User{
		ClerkID:   "clerk_test_1",
		Email:     ptr("test1@example.com"),
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
	}

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		// Insert test data
		_, err := tx.NewInsert().Model(testUser).Exec(ctx)
		require.NoError(s.T(), err)

		// Test getting the user
		user, err := services.GetUserByID(ctx, tx, testUser.ID)
		require.NoError(s.T(), err)
		require.Equal(s.T(), testUser.ID, user.ID)
		require.Equal(s.T(), "Test", user.FirstName)
		require.Equal(s.T(), "User", user.LastName)

		// test with non-existent user
		_, err = services.GetUserByID(ctx, tx, uuid.New())
		require.Error(s.T(), err)
		require.True(s.T(), errors.Is(err, sql.ErrNoRows))

		return nil
	})
	require.NoError(s.T(), err)
}

func (s *TestSuite) TestApplyActionToUser() {
	ctx := context.Background()

	testUser := &models.User{
		ClerkID:   "clerk_test_3",
		Email:     ptr("test3@example.com"),
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		Blocked:   false,
	}

	err := models.WithRollBackOnlyTransaction(ctx, s.application.DB, func(tx *bun.Tx) error {
		// Insert test data
		_, err := tx.NewInsert().Model(testUser).Exec(ctx)
		require.NoError(s.T(), err)

		// Test block action
		blockedUser, err := services.ApplyActionToUser(ctx, tx, testUser.ID, "block")
		require.NoError(s.T(), err)
		require.True(s.T(), blockedUser.Blocked)

		// Verify in database
		var dbUser models.User
		err = tx.NewSelect().Model(&dbUser).Where("id = ?", testUser.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.True(s.T(), dbUser.Blocked)

		// Test unblock action
		unblockedUser, err := services.ApplyActionToUser(ctx, tx, testUser.ID, "unblock")
		require.NoError(s.T(), err)
		require.False(s.T(), unblockedUser.Blocked)

		// Verify in database
		err = tx.NewSelect().Model(&dbUser).Where("id = ?", testUser.ID).Scan(ctx)
		require.NoError(s.T(), err)
		require.False(s.T(), dbUser.Blocked)

		// test invalid action
		_, err = services.ApplyActionToUser(ctx, tx, testUser.ID, "invalid-action")
		require.Error(s.T(), err)

		return nil
	})
	require.NoError(s.T(), err)
}

// convert a string to a pointer string
func ptr(s string) *string {
	return &s
}

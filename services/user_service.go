package services

import (
	"context"
	"rashikzaman/api/models"
	"strconv"
	"time"

	"github.com/uptrace/bun"
)

func CreateUserFromClerk(
	ctx context.Context, db bun.IDB, clerkID, email, firstName, lastName, birthday string,
) (*models.User, error) {
	user := &models.User{
		Email:     &email,
		ClerkID:   clerkID,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user",
	}

	if birthday != "" {
		millis, err := strconv.ParseInt(birthday, 10, 64)
		if err != nil {
			return nil, err
		}

		dateOfBirth := time.UnixMilli(millis)
		user.DateOfBirth = &dateOfBirth
	}

	err := models.Create(ctx, db, user)

	return user, err
}

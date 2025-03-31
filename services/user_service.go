package services

import (
	"context"
	"database/sql"
	"rashikzaman/api/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/pkg/errors"
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

func FetchUsersForAdmin(ctx context.Context, db bun.IDB, queryParam models.QueryParam) ([]models.User, int, error) {
	users := []models.User{}

	query := db.NewSelect().
		Model(&users)

	count, err := queryParam.Pagination.BuildPaginationQuery(ctx, query)
	if err != nil {
		return users, 0, errors.Wrap(err, err.Error())
	}

	if len(queryParam.Relations) != 0 {
		for _, relation := range queryParam.Relations {
			query.Relation(relation)
		}
	}

	err = query.Scan(ctx)
	if err != nil {
		return users, 0, errors.Wrap(err, err.Error())
	}

	return users, count, nil
}

func GetUserByID(ctx context.Context, db bun.IDB, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	user.ID = userID

	err := models.SelectByID(ctx, db, userID, user, models.QueryParam{
		Relations: []string{"UserLocations"},
		Alias:     "user"})

	return user, err
}

func ApplyActionToUser(ctx context.Context, db bun.IDB, userID uuid.UUID, action string) (*models.User, error) {
	user, err := GetUserByID(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	if action == "block" {
		user.Blocked = true
	} else if action == "unblock" {
		user.Blocked = false
	}

	err = models.Update(ctx, db, user)

	return user, err
}

func UpdateMe(ctx context.Context, db bun.IDB, existingUser *models.User, userBody models.User) (*models.User, error) {
	existingUser.FirstName = userBody.FirstName
	existingUser.LastName = userBody.LastName
	existingUser.PhoneNumber = userBody.PhoneNumber
	existingUser.ReceiveSMSNotification = userBody.ReceiveSMSNotification

	err := models.Update(ctx, db, existingUser)

	return existingUser, err
}

func CreateOrUpdateUserLocation(
	ctx context.Context, db bun.IDB, userLocationBody *models.UserLocation, userID uuid.UUID,
) (
	*models.UserLocation, error) {
	usrLocation := &models.UserLocation{}

	err := db.NewSelect().Model(usrLocation).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			userLocationBody.Location = models.PostgisGeometry{Geometry: orb.Point{userLocationBody.Longitude, userLocationBody.Latitude}, SRID: 4326}
			userLocationBody.UserID = userID
			err := models.Create(ctx, db, userLocationBody)

			return userLocationBody, err
		} else {
			return nil, err
		}
	}

	usrLocation.Latitude = userLocationBody.Latitude
	usrLocation.Longitude = userLocationBody.Longitude
	usrLocation.Location = models.PostgisGeometry{Geometry: orb.Point{userLocationBody.Longitude, userLocationBody.Latitude}, SRID: 4326}
	usrLocation.FormattedAddress = userLocationBody.FormattedAddress

	err = models.Update(ctx, db, usrLocation)

	return usrLocation, err
}

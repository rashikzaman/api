package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	Base
	ClerkID                string       `json:"clerk_id"`
	FirstName              string       `json:"first_name"`
	LastName               string       `json:"last_name"`
	Email                  *string      `json:"email" bun:"email,unique"`
	PhoneNumber            *string      `json:"phone_number" bun:"phone_number,unique"`
	DateOfBirth            *time.Time   `json:"date_of_birth"`
	Role                   string       `json:"role"`
	Blocked                bool         `json:"blocked"`
	UserLocations          UserLocation `bun:"rel:has-one,join:id=user_id" json:"user_location"`
	ReceiveSMSNotification bool         `json:"receive_sms_notification"`
}

type UserLocation struct {
	Base
	UserID           uuid.UUID       `bun:"type:uuid" json:"user_id"`
	User             *User           `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Latitude         float64         `json:"latitude"`
	Longitude        float64         `json:"longitude"`
	Location         PostgisGeometry `bun:"type:location" json:"-"`
	FormattedAddress string          `json:"formatted_address"`
}

func GetUserByClerkID(ctx context.Context, db bun.IDB, clerkID string) (*User, error) {
	user := &User{}

	err := db.NewSelect().Model(user).Where("clerk_id = ?", clerkID).Scan(ctx)
	return user, err

}

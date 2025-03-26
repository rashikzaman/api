package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	Base
	ClerkID     string     `json:"clerkId"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"Last_name"`
	Email       *string    `json:"email" bun:"email,unique"`
	PhoneNumber *string    `json:"phone_number" bun:"phone_number,unique"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Role        string     `json:"role"`
}

func GetUserByClerkID(ctx context.Context, db bun.IDB, clerkID string) (*User, error) {
	user := &User{}

	err := db.NewSelect().Model(user).Where("clerk_id = ?", clerkID).Scan(ctx)
	return user, err

}

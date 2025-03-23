package models

import (
	"time"
)

type User struct {
	Base
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"LastName"`
	Email       *string    `json:"email" bun:"email,unique"`
	PhoneNumber *string    `json:"phone_number" bun:"phone_number,unique"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Role        string     `json:"role"`
}

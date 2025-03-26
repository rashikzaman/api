package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rashikzaman/api/config"
	"rashikzaman/api/models"
	"rashikzaman/api/services"
	"time"

	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/uptrace/bun"
)

// ClerkUserWebhookEvent represents the top-level Clerk webhook event for user operations
type ClerkUserWebhookEvent struct {
	Data            UserData        `json:"data"`
	EventAttributes EventAttributes `json:"event_attributes"`
	Object          string          `json:"object"`
	Timestamp       int64           `json:"timestamp"`
	Type            string          `json:"type"`
}

// UserData represents a Clerk user object
type UserData struct {
	Birthday              string                 `json:"birthday"`
	CreatedAt             int64                  `json:"created_at"`
	EmailAddresses        []EmailAddress         `json:"email_addresses"`
	ExternalID            string                 `json:"external_id"`
	FirstName             string                 `json:"first_name"`
	Gender                string                 `json:"gender"`
	ID                    string                 `json:"id"`
	ImageURL              string                 `json:"image_url"`
	LastName              string                 `json:"last_name"`
	LastSignInAt          int64                  `json:"last_sign_in_at"`
	Object                string                 `json:"object"`
	PasswordEnabled       bool                   `json:"password_enabled"`
	PhoneNumbers          []interface{}          `json:"phone_numbers"`
	PrimaryEmailAddressID string                 `json:"primary_email_address_id"`
	PrimaryPhoneNumberID  *string                `json:"primary_phone_number_id"` // Using pointer for null values
	PrivateMetadata       map[string]interface{} `json:"private_metadata"`
	ProfileImageURL       string                 `json:"profile_image_url"`
	PublicMetadata        map[string]interface{} `json:"public_metadata"`
	TwoFactorEnabled      bool                   `json:"two_factor_enabled"`
	UnsafeMetadata        map[string]interface{} `json:"unsafe_metadata"`
	UpdatedAt             int64                  `json:"updated_at"`
	Username              *string                `json:"username"` // Using pointer for null values
}

type EmailAddress struct {
	EmailAddress string `json:"email_address"`
	ID           string `json:"id"`
	Object       string `json:"object"`
}

// EventAttributes contains additional information about the event
type EventAttributes struct {
	HTTPRequest HTTPRequestInfo `json:"http_request"`
}

// HTTPRequestInfo contains information about the HTTP request that triggered the event
type HTTPRequestInfo struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

func (e *ClerkUserWebhookEvent) GetEventTime() time.Time {
	return time.UnixMilli(e.Timestamp)
}

func (u *UserData) GetCreatedTime() time.Time {
	return time.UnixMilli(u.CreatedAt)
}

func (u *UserData) GetLastSignInTime() time.Time {
	return time.UnixMilli(u.LastSignInAt)
}

func (u *UserData) GetUpdatedTime() time.Time {
	return time.UnixMilli(u.UpdatedAt)
}

func (ac *Controller) ClerkWebHook(c *gin.Context) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// Create a new Svix webhook instance
	wh, err := svix.NewWebhook(config.GetClerkSigningSecretKey())
	if err != nil {
		log.Printf("Error creating Svix webhook: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	err = wh.Verify(body, c.Request.Header)
	if err != nil {
		log.Printf("Error: Could not verify webhook: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	var event ClerkUserWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if event.Type == "user.created" {
		user := event.Data

		err = models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
			_, err := services.CreateUserFromClerk(c, tx, user.ID, user.EmailAddresses[0].EmailAddress, user.FirstName, user.LastName, user.Birthday)
			if err != nil {
				log.Printf("Error creating user from clerk webhook: %v", err)
				return err
			}

			return nil
		})
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
}

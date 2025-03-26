package models

import "github.com/google/uuid"

type TaskMedia struct {
	Base
	MimeType string    `json:"mime_type"`
	Link     string    `json:"link"`
	TaskID   uuid.UUID `json:"task_id"`
	Base64   string    `bun:"-" json:"base64"`
}

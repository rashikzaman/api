package models

import "github.com/google/uuid"

type PostMedia struct {
	Base
	MimeType string
	Link     string
	PostID   uuid.UUID
}

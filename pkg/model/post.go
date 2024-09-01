package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id             uuid.UUID `json:"id"`
	OwnerId        uuid.UUID `json:"owner_id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	Content        string    `json:"content"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	Hidden         bool      `json:"hidden"`
	Deleted        bool      `json:"deleted"`
	RecoveryWindow time.Time `json:"recovery_window"`
}

func (p *Post) Validate() error {
	// TODO implement
	return nil
}

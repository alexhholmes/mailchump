package model

import "github.com/google/uuid"

type Post struct {
	Id           uuid.UUID   `json:"id"`
	NewsletterID uuid.UUID   `json:"newsletter_id"`
	AuthorIDs    []uuid.UUID `json:"author_ids"`
	Title        string      `json:"title"`
	Slug         string      `json:"slug"`
	Content      string      `json:"content"`
	Created      string      `json:"created"`
	Updated      string      `json:"updated"`
	Hidden       bool        `json:"hidden"`
}

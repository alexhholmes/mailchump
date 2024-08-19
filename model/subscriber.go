package model

import (
	"database/sql"
	"fmt"

	"mailchump/api"
)

// Subscriber defines model for Subscriber.
type Subscriber struct {
	Email string `json:"email"`
}

func (s *Subscriber) Validate() error {
	if s.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

// FromSubscription converts the APIs Subscription to a Subscriber.
func (s *Subscriber) FromSubscription(sub api.Subscription) {
	s.Email = sub.Email
}

// Read reads the subscriber from postgres `subscribers` table.
func (s *Subscriber) Read(db *sql.DB) error {
	row := db.QueryRow("SELECT email FROM subscribers WHERE email = $1", s.Email)
	if err := row.Scan(&s.Email); err != nil {
		return fmt.Errorf("failed to read subscriber: %w", err)
	}

	return nil
}

// Write writes the subscriber to postgres `subscribers` table.
func (s *Subscriber) Write(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO subscribers (email) VALUES ($1)", s.Email)
	if err != nil {
		return fmt.Errorf("failed to write subscriber: %w", err)
	}

	return nil
}

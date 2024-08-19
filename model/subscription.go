package model

import (
	"database/sql"
	"errors"
	"fmt"

	"mailchump/gen"
)

var (
	ErrCreateSubscriptionAlreadyExists = errors.New("cannot create subscription, already exists")
	ErrRemoveSubscriptionDoesNotExist  = errors.New("cannot remove subscription, does not exist")
	ErrRemoveSubscriptionForbidden     = errors.New("cannot remove subscription, forbidden")
)

// Subscription defines model for Subscription.
type Subscription struct {
	Email string `json:"email"`
}

func (s *Subscription) Validate() error {
	if s.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

// FromReq converts the APIs Subscription to a Subscription.
func (s *Subscription) FromReq(sub gen.SubscriptionRequest) {
	s.Email = sub.Email
}

// Create creates the subscriber in postgres `subscribers` table.
func (s *Subscription) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO subscribers (email) VALUES ($1)", s.Email)
	if err != nil {
		return fmt.Errorf("failed to write subscriber: %w", err)
	}

	return nil
}

func (s *Subscription) Remove(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM subscribers WHERE email = $1", s.Email)
	if err != nil {
		return fmt.Errorf("failed to remove subscriber: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrRemoveSubscriptionDoesNotExist
	}

	return nil
}

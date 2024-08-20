package model

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"mailchump/gen"
)

var (
	ErrCreateSubscriptionAlreadyExists = errors.New("cannot create subscription, already exists")
	ErrRemoveSubscriptionDoesNotExist  = errors.New("cannot remove subscription, does not exist")
	ErrRemoveSubscriptionForbidden     = errors.New("cannot remove subscription, forbidden")
)

// Subscription defines model for user's Subscription to a newsletter.
// An email can be used for multiple newsletters, but a newsletter can
// only have one subscription per email.
type Subscription struct {
	ID         uuid.UUID `json:"id"`
	Newsletter uuid.UUID `json:"newsletter"`
	Email      string    `json:"email"`
	Active     bool      `json:"active"`
	From       time.Time `json:"from"`
	Until      time.Time `json:"until"`
}

// Validate performs runtime checks on Subscription fields.
func (s *Subscription) Validate() error {
	if s.Email == "" {
		return fmt.Errorf("email is required")
	} else if _, err := mail.ParseAddress(s.Email); err != nil {
		return fmt.Errorf("email is invalid")
	}

	if s.Active {
		if time.Time.IsZero(s.From) {
			return fmt.Errorf("from is required")
		}
		if time.Time.IsZero(s.Until) {
			return fmt.Errorf("until is required")
		}
	} else {
		if time.Time.IsZero(s.From) {
			return fmt.Errorf("from is required")
		}
		if !time.Time.IsZero(s.Until) {
			return fmt.Errorf("until must be zero")
		}
	}

	return nil
}

// FromReq converts the APIs Subscription to a Subscription.
func (s *Subscription) FromReq(sub gen.SubscriptionRequest) {
	s.Email = sub.Email
}

// Create creates the subscriber in pgdb `subscribers` table.
func (s *Subscription) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO subscribers (email) VALUES ($1)", s.Email)
	if err != nil {
		return fmt.Errorf("failed to write subscriber: %w", err)
	}

	return nil
}

func (s *Subscription) Get(db *sql.DB) error {
	row := db.QueryRow("SELECT email FROM subscribers WHERE email = $1", s.Email)
	err := row.Scan(&s.Email)
	if err != nil {
		return fmt.Errorf("failed to get subscriber: %w", err)
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

package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"mailchump/pgdb"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNewsletterAlreadyExists = errors.New("newsletter already exists")
	ErrNewsletterNotFound      = errors.New("newsletter not found")
	ErrNoChanges               = errors.New("no changes made")
)

type Newsletter struct {
	ID             uuid.UUID   `json:"id"`
	OwnerID        uuid.UUID   `json:"owner_id"`
	AuthorIDs      []uuid.UUID `json:"author_ids"`
	Title          string      `json:"title"`
	Slug           string      `json:"slug"`
	Description    string      `json:"description"`
	Created        time.Time   `json:"created"`
	Updated        time.Time   `json:"updated"`
	PostCount      int         `json:"post_count"`
	Hidden         bool        `json:"hidden"`
	Deleted        bool        `json:"deleted"`
	RecoveryWindow time.Time   `json:"recovery_window"`
}

// Validate performs runtime checks on Newsletter fields.
func (n *Newsletter) Validate() error {
	// TODO
	return nil
}

func (n *Newsletter) Create(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	n.ID = uuid.New()
	n.Created = now
	n.Updated = now

	// Perform a transaction that inserts the newsletter and then adds the owner as an author
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer pgdb.HandleTxError(err, tx)

	// Create the newsletter
	res, err := tx.ExecContext(ctx,
		`INSERT INTO newsletters (id, owner, title, slug, description, created, updated, post_count, hidden)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING`,
		n.ID, n.OwnerID, n.Title, n.Slug, n.Description, n.Created, n.Updated, n.PostCount, n.Hidden,
	)
	if err != nil {
		return err
	}

	// Check if the newsletter was inserted
	if affected, err := res.RowsAffected(); err != nil {
		// DB driver does not support RowsAffected
		log.Fatal(err)
	} else if affected == 0 {
		return ErrNewsletterAlreadyExists
	}

	// Add the owner as an author of the newsletter
	if _, err = tx.ExecContext(ctx,
		`INSERT INTO newsletter_authors (newsletter, author)
		VALUES ($1, $2)`,
		n.ID, n.OwnerID,
	); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (n *Newsletter) Delete(ctx context.Context, db *sql.DB) error {
	// TODO make sure cascade delete works
	res, err := db.ExecContext(ctx,
		`UPDATE newsletters SET deleted = true WHERE id = $1`,
		n.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNewsletterNotFound
		}
		return err
	}

	// Notify caller that no changes were made, may want to give user a more specific
	// message/notification.
	if affected, err := res.RowsAffected(); err != nil {
		// DB driver does not support RowsAffected
		log.Fatal(err)
	} else if affected == 0 {
		return ErrNoChanges
	}

	return nil
}

func (n *Newsletter) Hide(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx,
		`UPDATE newsletters SET hidden = true WHERE id = $1`,
		n.ID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNewsletterNotFound
		}
		return err
	}

	return nil
}

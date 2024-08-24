package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"reflect"
	"time"

	"github.com/google/uuid"
	"mailchump/api/gen"
	"mailchump/api/util"
	"mailchump/pgdb"
)

type Newsletters []Newsletter

// ToResponse converts a slice of Newsletters to a slice of NewsletterResponse. The user
// parameter is used to determine if all fields should be shown (i.e. the user owns the
// newsletter).
func (n *Newsletters) ToResponse(user util.Key) *[]gen.NewsletterResponse {
	var resp []gen.NewsletterResponse
	for _, newsletter := range *n {
		resp = append(resp, newsletter.ToResponse(user))
	}
	return &resp
}

// GetAll fetches the entire newsletters table from the database.
func (n *Newsletters) GetAll(ctx context.Context, db *sql.DB) error {
	// Fetch entire newsletters table
	rows, err := db.QueryContext(ctx,
		`SELECT 
		FROM newsletters
		WHERE deleted = false`,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	defer pgdb.HandleCloseResult(rows)

	for rows.Next() {
		newsletter := Newsletter{}

		s := reflect.ValueOf(&newsletter).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err := rows.Scan(columns...)
		if err != nil {
			log.Fatal(err)
		}

		*n = append(*n, newsletter)
	}

	// Check for errors during rows.Next
	if err = rows.Close(); err != nil {
		slog.Warn("failed to close rows", "error", err)
		return err
	}

	return nil
}

type Newsletter struct {
	Id             uuid.UUID   `json:"id"`
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
	// TODO implement
	return nil
}

func (n *Newsletter) ToResponse(user util.Key) gen.NewsletterResponse {
	// Hide fields if the user is not an owner
	var (
		owner   *string
		hidden  *bool
		deleted *bool
	)
	if user.String() == n.OwnerID.String() {
		s := n.OwnerID.String()
		owner = &s
		hidden = &n.Hidden
		deleted = &n.Deleted
	}

	return gen.NewsletterResponse{
		Authors: func() []string {
			var authors []string
			for _, a := range n.AuthorIDs {
				authors = append(authors, a.String())
			}
			return authors
		}(),
		CreatedAt:   n.Created.String(),
		Deleted:     deleted,
		Description: n.Description,
		Hidden:      hidden,
		Id:          n.Id.String(),
		Owner:       owner,
		PostCount:   n.PostCount,
		Slug:        n.Slug,
		Title:       n.Title,
		UpdatedAt:   n.Updated.String(),
	}
}

func (n *Newsletter) Get(ctx context.Context, db *sql.DB) error {
	// Fetch the newsletter from the database
	err := db.QueryRowContext(ctx,
		`SELECT owner, title, slug, description, created, updated, post_count, hidden, deleted, recovery_window
		FROM newsletters
		WHERE id = $1`,
		n.Id,
	).Scan(&n.OwnerID, &n.Title, &n.Slug, &n.Description, &n.Created, &n.Updated, &n.PostCount, &n.Hidden, &n.Deleted, &n.RecoveryWindow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	// Fetch all authors of the newsletter
	rows, err := db.QueryContext(ctx,
		`SELECT author
		FROM newsletter_authors
		WHERE newsletter = $1`,
		n.Id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	defer pgdb.HandleCloseResult(rows)

	for rows.Next() {
		var author uuid.UUID
		if err = rows.Scan(&author); err != nil {
			return err
		}
		n.AuthorIDs = append(n.AuthorIDs, author)
	}

	return nil
}

func (n *Newsletter) GetOwnerID(ctx context.Context, db *sql.DB) error {
	// Fetch the owner of the newsletter
	err := db.QueryRowContext(ctx,
		`SELECT owner
		FROM newsletters
		WHERE id = $1`,
		n.Id,
	).Scan(&n.OwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (n *Newsletter) Create(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()

	n.Id = uuid.New()
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
		n.Id, n.OwnerID, n.Title, n.Slug, n.Description, n.Created, n.Updated, n.PostCount, n.Hidden,
	)
	if err != nil {
		return err
	}

	// Check if the newsletter was inserted
	if affected, err := res.RowsAffected(); err != nil {
		// DB driver does not support RowsAffected
		log.Fatal(err)
	} else if affected == 0 {
		return ErrAlreadyExists
	}

	// Add the owner as an author of the newsletter
	if _, err = tx.ExecContext(ctx,
		`INSERT INTO newsletter_authors (newsletter, author)
		VALUES ($1, $2)`,
		n.Id, n.OwnerID,
	); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (n *Newsletter) Delete(ctx context.Context, db *sql.DB) error {
	// Delete the newsletter from the database
	res, err := db.ExecContext(ctx,
		`UPDATE newsletters
		SET deleted = true, recovery_window = $2
		WHERE id = $1
		AND deleted = false`,
		n.Id, time.Now().AddDate(0, 0, 7), // TODO add recovery time to config
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	} else if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return ErrNotFound
	}

	// Check if the newsletter was deleted
	if affected, err := res.RowsAffected(); err != nil {
		// DB driver does not support RowsAffected
		log.Fatal(err)
	} else if affected == 0 {
		return ErrNotFound
	}

	return nil
}

// Hide changes the hidden field of a Newsletter to the opposite of its current value.
func (n *Newsletter) Hide(ctx context.Context, db *sql.DB) error {
	if n.Id == uuid.Nil {
		return errors.New("nil UUID")
	}

	// Flip the hidden field of the newsletter
	res, err := db.QueryContext(ctx,
		`UPDATE newsletters
		SET hidden = NOT hidden
		WHERE id = $1
		RETURNING hidden`,
		n.Id,
	)
	if err != nil {
		return err
	}

	// Get newsletter hidden status from result
	if res.Next() {
		if err = res.Scan(&n.Hidden); err != nil {
			return err
		}

		return nil
	}

	return ErrNotFound
}

func (n *Newsletter) IsOwner(ctx context.Context, db *sql.DB) (bool, error) {
	// Check that the user is the newsletter owner
	err := n.GetOwnerID(ctx, db)
	if err != nil {
		return false, err
	}

	user := ctx.Value(util.ContextUser).(string)
	return user == n.OwnerID.String(), nil
}

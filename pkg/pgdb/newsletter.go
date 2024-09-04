package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"mailchump/pkg/model"
)

// GetAllNewsletters fetches the entire newsletters table from the database.
func (c *Client) GetAllNewsletters(ctx context.Context) (model.Newsletters, error) {
	// Fetch entire newsletters table
	rows, err := c.db.QueryContext(ctx,
		`SELECT 
		FROM newsletters
		WHERE deleted = false`,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	defer HandleCloseResult(rows)

	var n []model.Newsletter
	for rows.Next() {
		newsletter, _ := MapStruct[model.Newsletter](rows)
		n = append(n, newsletter)
	}

	return n, nil
}

func (c *Client) GetNewsletterById(
	ctx context.Context,
	id string,
) (model.Newsletter, error) {
	// Fetch the newsletter from the database
	row := c.db.QueryRowContext(ctx,
		`SELECT *
		FROM newsletters
		WHERE id = $1`,
		id,
	)

	newsletter, err := MapStruct[model.Newsletter](row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Newsletter{}, model.ErrNotFound
		}
		return model.Newsletter{}, err
	}

	// TODO get all authors from the posts of this newsletter
	newsletter.AuthorIDs = []uuid.UUID{newsletter.OwnerID}

	return newsletter, nil
}

func (c *Client) GetNewsletterOwnerID(ctx context.Context, id string) (uuid.UUID, error) {
	var ownerID uuid.UUID

	// Fetch the owner of the newsletter
	err := c.db.QueryRowContext(ctx,
		`SELECT owner
		FROM newsletters
		WHERE id = $1`,
		id,
	).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ownerID, model.ErrNotFound
		}
		return ownerID, err
	}

	return ownerID, nil
}

func (c *Client) CreateNewsletter(ctx context.Context, n model.Newsletter) (
	model.Newsletter,
	error,
) {
	now := time.Now().UTC()

	n.Id = uuid.New()
	n.Created = now
	n.Updated = now

	// Perform a transaction that inserts the newsletter and then adds the owner as an author
	tx, err := c.db.Begin()
	if err != nil {
		return model.Newsletter{}, err
	}
	defer HandleTxError(err, tx)

	// Create the newsletter
	res, err := tx.ExecContext(ctx,
		`INSERT INTO newsletters (id, owner, title, slug, description, created, updated, post_count, hidden)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING`,
		n.Id, n.OwnerID, n.Title, n.Slug, n.Description, n.Created, n.Updated, n.PostCount, n.Hidden,
	)
	if err != nil {
		return model.Newsletter{}, err
	}

	// Check if the newsletter was inserted
	if affected, err := res.RowsAffected(); err != nil {
		// db driver does not support RowsAffected
		log.Fatal(err)
	} else if affected == 0 {
		return model.Newsletter{}, model.ErrAlreadyExists
	}

	// Add the owner as an author of the newsletter
	if _, err = tx.ExecContext(ctx,
		`INSERT INTO newsletter_authors (newsletter, author)
		VALUES ($1, $2)`,
		n.Id, n.OwnerID,
	); err != nil {
		return model.Newsletter{}, err
	}

	if err = tx.Commit(); err != nil {
		return model.Newsletter{}, err
	}

	return n, nil
}

func (c *Client) DeleteNewsletter(ctx context.Context, id string) error {
	// Delete the newsletter from the database
	res, err := c.db.ExecContext(ctx,
		`UPDATE newsletters
		SET deleted = true, recovery_window = $2
		WHERE id = $1
		AND deleted = false`,
		id, time.Now().AddDate(0, 0, 7), // TODO add recovery time to config
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrNotFound
		}
		return err
	} else if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return model.ErrNotFound
	}

	// Check if the newsletter was deleted
	if affected, err := res.RowsAffected(); err != nil {
		// db driver does not support RowsAffected
		log.Fatal(err)
	} else if affected == 0 {
		return model.ErrNotFound
	}

	return nil
}

// HideNewsletter changes the hidden field of a Newsletter to the opposite of
// its current value.
func (c *Client) HideNewsletter(ctx context.Context, id string) (
	isHidden bool,
	err error,
) {
	// Flip the hidden field of the newsletter
	res, err := c.db.QueryContext(ctx,
		`UPDATE newsletters
		SET hidden = NOT hidden
		WHERE id = $1
		RETURNING hidden`,
		id,
	)
	if err != nil {
		return false, err
	}
	defer HandleCloseResult(res)

	// Get newsletter hidden status from result
	if res.Next() {
		if err = res.Scan(&isHidden); err != nil {
			return false, err
		}
		return isHidden, nil
	}

	return false, model.ErrNotFound
}

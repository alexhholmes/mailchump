package newsletters

import "errors"

var (
	ErrNewsletterAlreadyExists = errors.New("newsletter already exists")
	ErrNewsletterNotFound      = errors.New("newsletter not found")
	ErrNoChanges               = errors.New("no changes made")
)

package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/util"
)

type Newsletter struct {
	Id             uuid.UUID   `json:"id" table:"newsletters"`
	OwnerID        uuid.UUID   `json:"owner_id" table:"newsletters"`
	AuthorIDs      []uuid.UUID `json:"author_ids" table:"authors"`
	Title          string      `json:"title" table:"newsletters"`
	Slug           string      `json:"slug" table:"newsletters"`
	Description    string      `json:"description" table:"newsletters"`
	Created        time.Time   `json:"created" table:"newsletters"`
	Updated        time.Time   `json:"updated" table:"newsletters"`
	PostCount      int         `json:"post_count" table:"newsletters"`
	Hidden         bool        `json:"hidden" table:"newsletters"`
	Deleted        bool        `json:"deleted" table:"newsletters"`
	RecoveryWindow time.Time   `json:"recovery_window" table:"newsletters"`
}

// Validate performs runtime checks on Newsletter fields.
func (n *Newsletter) Validate() error {
	// TODO implement
	return nil
}

// ToResponse converts a Newsletter to a gen.NewsletterResponse. This will hide the
// following fields if the user is not the owner: Hidden and Deleted.
func (n *Newsletter) ToResponse(user string) gen.NewsletterResponse {
	resp := gen.NewsletterResponse{
		Authors: func() []string {
			var authors []string
			for _, a := range n.AuthorIDs {
				authors = append(authors, a.String())
			}
			return authors
		}(),
		CreatedAt:   n.Created.String(),
		Deleted:     nil,
		Description: n.Description,
		Hidden:      nil,
		Id:          n.Id.String(),
		Owner:       n.OwnerID.String(),
		PostCount:   n.PostCount,
		Slug:        n.Slug,
		Title:       n.Title,
		UpdatedAt:   n.Updated.String(),
	}
	// Hide fields if the user is not an owner
	if user == n.OwnerID.String() {
		resp.Hidden = &n.Hidden
		resp.Deleted = &n.Deleted
	}

	return resp
}

func (n *Newsletter) IsOwner(ctx context.Context) (bool, error) {
	// Check that the user is the newsletter owner
	user := ctx.Value(util.ContextUser).(string)
	return user == n.OwnerID.String(), nil
}

type Newsletters []Newsletter

// ToResponse converts a slice of Newsletters to a slice of NewsletterResponse. The user
// parameter is used to determine if all fields should be shown (i.e. the user owns the
// newsletter).
func (n *Newsletters) ToResponse(user string) []gen.NewsletterResponse {
	resp := make([]gen.NewsletterResponse, 0, len(*n))
	for _, newsletter := range *n {
		resp = append(resp, newsletter.ToResponse(user))
	}
	return resp
}

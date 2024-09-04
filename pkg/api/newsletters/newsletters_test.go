package newsletters

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/util"
	"mailchump/pkg/model"
)

type NewslettersTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestNewsletterTestSuite(t *testing.T) {
	suite.Run(t, new(NewslettersTestSuite))
}

func (s *NewslettersTestSuite) SetupTest() {
	s.ctx = context.WithValue(
		context.Background(),
		util.ContextUser,
		"00000000-0000-0000-0000-000000000000",
	)
	s.ctx = context.WithValue(
		s.ctx,
		util.ContextLogger,
		slog.Default(),
	)
}

func (s *NewslettersTestSuite) TestGetAllNewsletters() {
	t := s.T()

	exp := gen.AllNewsletterResponse{
		// We want an empty array not a nil value
		Newsletters: []gen.NewsletterResponse{},
		Count:       0,
	}

	m := NewMockNewsletterStore(t)
	m.EXPECT().GetAllNewsletters(s.ctx).Return(model.Newsletters{}, nil)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/newsletters", nil).WithContext(s.ctx)

	h := NewsletterHandler{DB: m}
	h.GetNewsletters(resp, req)

	var out gen.AllNewsletterResponse
	assert.Equal(t, 200, resp.Code)
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	assert.Equal(t, exp, out)
}

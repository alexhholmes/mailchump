package newsletters

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/util"
	"mailchump/pkg/model"
)

// Vars for nullable boolean values
var (
	True  = true
	False = false
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

	testcases := []struct {
		name string
		exp  gen.AllNewsletterResponse
		mock func(*MockNewsletterStore) *MockNewsletterStore
	}{
		{
			name: "empty",
			exp: gen.AllNewsletterResponse{
				Newsletters: []gen.NewsletterResponse{},
				Count:       0,
			},
			mock: func(m *MockNewsletterStore) *MockNewsletterStore {
				m.EXPECT().GetAllNewsletters(s.ctx).Return(model.Newsletters{}, nil)
				return m
			},
		},
		{
			name: "single",
			exp: gen.AllNewsletterResponse{
				Newsletters: []gen.NewsletterResponse{
					{
						Id:        "00000000-0000-0000-0000-000000000000",
						Owner:     "00000000-0000-0000-0000-000000000000",
						CreatedAt: time.Time{}.String(),
						UpdatedAt: time.Time{}.String(),
						Deleted:   &False,
						Hidden:    &False,
					},
				},
				Count: 1,
			},
			mock: func(m *MockNewsletterStore) *MockNewsletterStore {
				m.EXPECT().GetAllNewsletters(s.ctx).
					Return(model.Newsletters{
						{
							Id:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
							OwnerID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						},
					}, nil)
				return m
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/newsletters", nil).WithContext(s.ctx)

			m := tc.mock(NewMockNewsletterStore(s.T()))
			h := NewsletterHandler{DB: m}
			h.GetNewsletters(resp, req)

			var out gen.AllNewsletterResponse

			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
			assert.Equal(t, 200, resp.Code)
			assert.Equal(t, tc.exp, out)
		})
	}
}

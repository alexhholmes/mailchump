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
	ctx         context.Context
	newsletters []model.Newsletter
}

func TestNewsletterTestSuite(t *testing.T) {
	suite.Run(t, func() *NewslettersTestSuite {
		s := &NewslettersTestSuite{
			newsletters: make([]model.Newsletter, 0, 1000),
		}

		for i := 0; i < cap(s.newsletters); i++ {
			if i%2 == 0 {
				s.newsletters = append(s.newsletters, model.Newsletter{
					Id:      uuid.New(),
					OwnerID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				})
			} else {
				s.newsletters = append(s.newsletters, model.Newsletter{
					Id:      uuid.New(),
					OwnerID: uuid.New(),
				})
			}
		}

		return s
	}())
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
						Id:        s.newsletters[0].Id.String(),
						Owner:     s.newsletters[0].OwnerID.String(),
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
					Return(s.newsletters[0:1], nil)
				return m
			},
		},
		{
			name: "multiple",
			exp: gen.AllNewsletterResponse{
				Newsletters: []gen.NewsletterResponse{
					{
						Id:        s.newsletters[0].Id.String(),
						Owner:     s.newsletters[0].OwnerID.String(),
						CreatedAt: time.Time{}.String(),
						UpdatedAt: time.Time{}.String(),
						Deleted:   &False,
						Hidden:    &False,
					},
					{
						Id:        s.newsletters[1].Id.String(),
						Owner:     s.newsletters[1].OwnerID.String(),
						CreatedAt: time.Time{}.String(),
						UpdatedAt: time.Time{}.String(),
						Deleted:   nil,
						Hidden:    nil,
					},
				},
				Count: 2,
			},
			mock: func(m *MockNewsletterStore) *MockNewsletterStore {
				m.EXPECT().GetAllNewsletters(s.ctx).
					Return(s.newsletters[0:2], nil)
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

package newsletters

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"mailchump/pkg/api/util"
	"mailchump/pkg/model"
)

func TestGetAllNewsletters(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		util.ContextUser,
		"00000000-0000-0000-0000-000000000000",
	)

	m := NewMockNewsletterStore(t)
	m.EXPECT().GetAllNewsletters(ctx).Return(model.Newsletters{}, nil)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/newsletters", nil).WithContext(ctx)

	h := NewsletterHandler{DB: m}
	h.GetNewsletters(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, `{"newsletters":[],"count":0}`, resp.Body.String())
}

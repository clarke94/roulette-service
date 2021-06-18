package table

import (
	"bytes"
	"context"
	"errors"
	"github.com/clarke94/roulette-service/internal/pkg/table"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		want       Handler
	}{
		{
			name:       "expect Handler to init",
			controller: mockController{},
			want: Handler{
				Controller: mockController{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.controller)

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

type mockController struct {
	GivenID    uuid.UUID
	GivenError error
}

func (m mockController) Create(_ context.Context, _ table.Table) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		body       []byte
		wantCode   int
	}{
		{
			name: "expect 201 given table created",
			controller: mockController{
				GivenID: uuid.New(),
			},
			body:     []byte(`{}`),
			wantCode: http.StatusCreated,
		},
		{
			name: "expect 400 given no body",
			controller: mockController{
				GivenID: uuid.New(),
			},
			body:     nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			body:     nil,
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = r

			h := NewHandler(tt.controller)
			h.Create(ctx)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

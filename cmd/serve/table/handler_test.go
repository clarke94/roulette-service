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
			body:     []byte(`{}`),
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

func TestHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		wantCode   int
	}{
		{
			name: "expect 200 given table created",
			controller: mockController{
				GivenList: []table.Table{
					{
						ID:         uuid.Nil,
						Name:       "Table 1",
						MaximumBet: 10000,
						MinimumBet: 1000,
						Currency:   "GBP",
					},
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = r

			h := NewHandler(tt.controller)
			h.List(ctx)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		body       []byte
		wantCode   int
	}{
		{
			name: "expect 200 given table updated",
			controller: mockController{
				GivenID: uuid.New(),
			},
			body:     []byte(`{}`),
			wantCode: http.StatusOK,
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
			body:     []byte(`{}`),
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
			h.Update(ctx)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		id         string
		wantCode   int
	}{
		{
			name: "expect 200 given table deleted",
			controller: mockController{
				GivenID: uuid.New(),
			},
			id:       uuid.New().String(),
			wantCode: http.StatusOK,
		},
		{
			name:       "expect 400 given invalid ID",
			controller: mockController{},
			id:         "foo",
			wantCode:   http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			id:       uuid.New().String(),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.controller)

			r := httptest.NewRequest(http.MethodDelete, "/"+tt.id, nil)
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			ctx.Request = r

			router.Handle(http.MethodDelete, "/:id", h.Delete)
			router.ServeHTTP(w, r)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

type mockController struct {
	GivenList  []table.Table
	GivenID    uuid.UUID
	GivenError error
}

func (m mockController) Delete(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

func (m mockController) Update(_ context.Context, _ table.Table) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

func (m mockController) List(_ context.Context) ([]table.Table, error) {
	return m.GivenList, m.GivenError
}

func (m mockController) Create(_ context.Context, _ table.Table) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

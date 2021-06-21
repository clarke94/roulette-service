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
				GivenID: uuid.New().String(),
			},
			body:     []byte(`{"name":"foo", "maximumBet": 1000, "minimumBet": 100, "currency":"GBP"}`),
			wantCode: http.StatusCreated,
		},
		{
			name: "expect 400 given no body",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			body:     nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			body:     []byte(`{"name":"foo", "maximumBet": 1000, "minimumBet": 100, "currency":"GBP"}`),
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
						ID:         "",
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
				GivenID: uuid.New().String(),
			},
			body:     []byte(`{"id":"42bb1490-d28e-11eb-b8bc-0242ac130003", "name":"foo", "maximumBet": 1000, "minimumBet": 100, "currency":"GBP"}`),
			wantCode: http.StatusOK,
		},
		{
			name: "expect 400 given no body",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			body:     nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			body:     []byte(`{"id":"42bb1490-d28e-11eb-b8bc-0242ac130003", "name":"foo", "maximumBet": 1000, "minimumBet": 100, "currency":"GBP"}`),
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
				GivenID: uuid.New().String(),
			},
			id:       "84b10ade-d28a-11eb-b8bc-0242ac130003",
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
			id:       "84b10ade-d28a-11eb-b8bc-0242ac130003",
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

			router.Handle(http.MethodDelete, "/:table", h.Delete)
			router.ServeHTTP(w, r)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

type mockController struct {
	GivenList  []table.Table
	GivenID    string
	GivenError error
}

func (m mockController) Delete(_ context.Context, _ string) (string, error) {
	return m.GivenID, m.GivenError
}

func (m mockController) Update(_ context.Context, _ table.Table) (string, error) {
	return m.GivenID, m.GivenError
}

func (m mockController) List(_ context.Context) ([]table.Table, error) {
	return m.GivenList, m.GivenError
}

func (m mockController) Create(_ context.Context, _ table.Table) (string, error) {
	return m.GivenID, m.GivenError
}

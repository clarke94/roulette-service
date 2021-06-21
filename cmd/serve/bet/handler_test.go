package bet

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clarke94/roulette-service/internal/pkg/bet"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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
		tableId    string
		body       []byte
		wantCode   int
	}{
		{
			name: "expect 201 given bet created",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			tableId:  uuid.New().String(),
			body:     []byte(`{"bet":"foo", "type":"foo", "amount": 10, "currency": "GBP"}`),
			wantCode: http.StatusCreated,
		},
		{
			name: "expect 400 given invalid table request",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			tableId:  uuid.New().String(),
			body:     []byte(`{}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given invalid table ID",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			tableId:  "foo",
			body:     []byte(`{"bet":"foo", "type":"foo", "amount": 10, "currency": "GBP"}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			tableId:  uuid.New().String(),
			body:     []byte(`{"bet":"foo", "type":"foo", "amount": 10, "currency": "GBP"}`),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.controller)

			r := httptest.NewRequest(http.MethodPost, "/"+tt.tableId, bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			ctx.Request = r

			router.Handle(http.MethodPost, "/:table", h.Create)
			router.ServeHTTP(w, r)

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
		tableId    string
		wantCode   int
	}{
		{
			name: "expect 200 given bet created",
			controller: mockController{
				GivenList: []bet.Bet{
					{
						ID:       "",
						TableID:  "",
						Amount:   10,
						Type:     bet.TypeStraight,
						Currency: "GBP",
					},
				},
			},
			tableId:  uuid.New().String(),
			wantCode: http.StatusOK,
		},
		{
			name: "expect 400 given invalid table ID",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			tableId:  "foo",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			tableId:  uuid.New().String(),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.controller)

			r := httptest.NewRequest(http.MethodGet, "/"+tt.tableId, nil)
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			ctx.Request = r

			router.Handle(http.MethodGet, "/:table", h.List)
			router.ServeHTTP(w, r)

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
		tableId    string
		body       []byte
		wantCode   int
	}{
		{
			name: "expect 200 given bet updated",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			tableId:  uuid.New().String(),
			body:     []byte(`{"id":"42bb1490-d28e-11eb-b8bc-0242ac130003", "bet":"foo", "type":"foo", "amount": 10, "currency": "GBP"}`),
			wantCode: http.StatusOK,
		},
		{
			name: "expect 400 given invalid table ID",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			tableId:  "foo",
			body:     []byte(`{"id":"42bb1490-d28e-11eb-b8bc-0242ac130003", "bet":"foo", "type":"foo", "amount": 10, "currency": "GBP"}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given invalid request",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			tableId:  uuid.New().String(),
			body:     []byte(`{}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			tableId:  uuid.New().String(),
			body:     []byte((`{"id":"42bb1490-d28e-11eb-b8bc-0242ac130003", "bet":"foo", "type":"foo", "amount": 10, "currency": "GBP"}`)),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.controller)

			r := httptest.NewRequest(http.MethodPut, "/"+tt.tableId, bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			ctx.Request = r

			router.Handle(http.MethodPut, "/:table", h.Update)
			router.ServeHTTP(w, r)

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
		tableId    string
		wantCode   int
	}{
		{
			name: "expect 200 given bet deleted",
			controller: mockController{
				GivenID: uuid.New().String(),
			},
			id:       uuid.New().String(),
			tableId:  uuid.New().String(),
			wantCode: http.StatusOK,
		},
		{
			name:       "expect 400 given invalid ID",
			controller: mockController{},
			id:         "foo",
			tableId:    uuid.New().String(),
			wantCode:   http.StatusBadRequest,
		},
		{
			name:       "expect 400 given invalid table ID",
			controller: mockController{},
			tableId:    "foo",
			id:         uuid.New().String(),
			wantCode:   http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			id:       uuid.New().String(),
			tableId:  uuid.New().String(),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.controller)

			r := httptest.NewRequest(http.MethodDelete, "/"+tt.tableId+"/"+tt.id, nil)
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			ctx.Request = r

			router.Handle(http.MethodDelete, "/:table/:bet", h.Delete)
			router.ServeHTTP(w, r)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_Play(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		tableId    string
		wantCode   int
	}{
		{
			name: "expect 200 given game played",
			controller: mockController{
				GivenResult: bet.Result{
					Number: 10,
					Color:  "red",
					Winners: []bet.Winner{
						{
							BetID:    uuid.New().String(),
							Amount:   10,
							Currency: "GBP",
						},
					},
				},
			},
			tableId:  uuid.New().String(),
			wantCode: http.StatusOK,
		},
		{
			name:       "expect 400 given invalid table ID",
			controller: mockController{},
			tableId:    "foo",
			wantCode:   http.StatusBadRequest,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenError: errors.New("foo"),
			},
			tableId:  uuid.New().String(),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.controller)

			r := httptest.NewRequest(http.MethodPost, "/"+tt.tableId+"/play", nil)
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			ctx.Request = r

			router.Handle(http.MethodPost, "/:table/play", h.Play)
			router.ServeHTTP(w, r)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

type mockController struct {
	GivenResult bet.Result
	GivenList   []bet.Bet
	GivenID     string
	GivenError  error
}

func (m mockController) Play(_ context.Context, _ string) (bet.Result, error) {
	return m.GivenResult, m.GivenError
}

func (m mockController) Delete(_ context.Context, _, _ string) (string, error) {
	return m.GivenID, m.GivenError
}

func (m mockController) Update(_ context.Context, _ bet.Bet) (string, error) {
	return m.GivenID, m.GivenError
}

func (m mockController) List(_ context.Context, _ string) ([]bet.Bet, error) {
	return m.GivenList, m.GivenError
}

func (m mockController) Create(_ context.Context, _ bet.Bet) (string, error) {
	return m.GivenID, m.GivenError
}

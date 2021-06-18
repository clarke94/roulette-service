package openapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

func TestHandler_Docs(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		wantCode   int
	}{
		{
			name: "expect 200 success",
			controller: mockController{
				GivenBytes: []byte(`<html></html>`),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenBytes: []byte(`<html></html>`),
				GivenError: errors.New("foo"),
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			h := NewHandler(tt.controller)

			h.Docs(ctx)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_Specification(t *testing.T) {
	tests := []struct {
		name       string
		controller ControllerProvider
		wantCode   int
	}{
		{
			name: "expect 200 success",
			controller: mockController{
				GivenBytes: []byte(`{"foo": "bar"}`),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "expect 400 given Controller error",
			controller: mockController{
				GivenBytes: []byte(`{"foo": "bar"}`),
				GivenError: errors.New("foo"),
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "expect 400 given no data found",
			controller: mockController{
				GivenBytes: nil,
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			h := NewHandler(tt.controller)

			h.Specification(ctx)

			if !cmp.Equal(w.Code, tt.wantCode) {
				t.Error(w.Code, tt.wantCode)
			}
		})
	}
}

type mockController struct {
	GivenBytes []byte
	GivenError error
}

func (m mockController) Specification() ([]byte, error) {
	return m.GivenBytes, m.GivenError
}

func (m mockController) Docs() ([]byte, error) {
	return m.GivenBytes, m.GivenError
}

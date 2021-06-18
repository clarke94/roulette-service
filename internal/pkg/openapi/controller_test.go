package openapi

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Controller
	}{
		{
			name: "expect Controller to init",
			want: Controller{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(logrus.New())
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(Controller{}), cmpopts.IgnoreFields(Controller{}, "Logger")) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(Controller{}), cmpopts.IgnoreFields(Controller{}, "Logger")))
			}
		})
	}
}

func TestController_Docs(t *testing.T) {
	tests := []struct {
		name       string
		controller Controller
		want       []byte
		wantErr    error
	}{
		{
			name: "expect success given valid doc",
			controller: Controller{
				docs: "docs/redoc_test.html",
			},
			want:    []byte(`<html lang="en-GB">foo</html>`),
			wantErr: nil,
		},
		{
			name: "expect fail given doc doesnt exist",
			controller: Controller{
				Logger: logrus.New(),
				docs:   "docs/foo",
			},
			want:    nil,
			wantErr: ErrDocumentation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.controller.Docs()

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestController_Specification(t *testing.T) {
	tests := []struct {
		name       string
		controller Controller
		want       []byte
		wantErr    error
	}{
		{
			name: "expect success given valid doc",
			controller: Controller{
				Logger:  logrus.New(),
				swagger: "docs/swagger_test.json",
			},
			want:    []byte(`{"foo": "bar"}`),
			wantErr: nil,
		},
		{
			name: "expect fail given spec doesnt exist",
			controller: Controller{
				Logger: logrus.New(),
				docs:   "docs/foo",
			},
			want:    nil,
			wantErr: ErrSpecification,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.controller.Specification()

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

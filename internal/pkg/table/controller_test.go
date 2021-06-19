package table

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
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
			want: Controller{
				Storage: mockStorage{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(logrus.New(), mockStorage{}, validator.New())
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(Controller{}), cmpopts.IgnoreFields(Controller{}, "Logger", "Validator")) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(Controller{}), cmpopts.IgnoreFields(Controller{}, "Logger", "Validator")))
			}
		})
	}
}

func TestController_Create(t *testing.T) {
	tests := []struct {
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		model     Table
		wantErr   error
	}{
		{
			name:      "expect success given valid table",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect fail given invalid maximum bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				MaximumBet: 1,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given invalid minimum bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				MaximumBet: 10,
				MinimumBet: 1,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given maximum bet less than minimum bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				MaximumBet: 10,
				MinimumBet: 20,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given invalid currency code",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				MaximumBet: 10,
				MinimumBet: 10,
				Currency:   "foo",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given storage error",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			model: Table{
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: ErrCreate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)
			_, err := c.Create(context.Background(), tt.model)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_List(t *testing.T) {
	tests := []struct {
		name       string
		Logger     *logrus.Logger
		Validator  *validator.Validate
		Storage    StorageProvider
		wantTables []Table
		wantErr    error
	}{
		{
			name:      "expect success given no tables found",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList: []Table{},
			},
			wantTables: []Table{},
			wantErr:    nil,
		},
		{
			name:      "expect success given tables found",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList: []Table{
					{
						Name:       "foo",
						MaximumBet: 10,
						MinimumBet: 10,
						Currency:   "GBP",
					},
				},
			},
			wantTables: []Table{
				{
					Name:       "foo",
					MaximumBet: 10,
					MinimumBet: 10,
					Currency:   "GBP",
				},
			},
			wantErr: nil,
		},
		{
			name:      "expect fail given storage error",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
				GivenList:  []Table{},
			},
			wantTables: []Table{},
			wantErr:    ErrList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)
			tables, err := c.List(context.Background())

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(tables, tt.wantTables) {
				t.Error(cmp.Diff(tables, tt.wantTables))
			}
		})
	}
}

func TestController_Update(t *testing.T) {
	tests := []struct {
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		model     Table
		wantErr   error
	}{
		{
			name:      "expect success given valid table",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				ID:         uuid.New(),
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect fail given no ID",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				Name:       "foo",
				MaximumBet: 1,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given invalid maximum bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				ID:         uuid.New(),
				Name:       "foo",
				MaximumBet: 1,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given invalid minimum bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				ID:         uuid.New(),
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 1,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given maximum bet less than minimum bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				ID:         uuid.New(),
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 20,
				Currency:   "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given invalid currency code",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Table{
				ID:         uuid.New(),
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 10,
				Currency:   "foo",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given storage error",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			model: Table{
				ID:         uuid.New(),
				Name:       "foo",
				MaximumBet: 10,
				MinimumBet: 10,
				Currency:   "GBP",
			},
			wantErr: ErrUpdate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)
			_, err := c.Update(context.Background(), tt.model)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockStorage struct {
	GivenList  []Table
	GivenID    uuid.UUID
	GivenError error
}

func (m mockStorage) Update(_ context.Context, _ Table) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

func (m mockStorage) List(_ context.Context) ([]Table, error) {
	return m.GivenList, m.GivenError
}

func (m mockStorage) Create(_ context.Context, _ Table) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

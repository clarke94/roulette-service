package bet

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
			got := New(logrus.New(), mockStorage{})
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(Controller{}), cmpopts.IgnoreFields(Controller{}, "Logger")) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(Controller{}), cmpopts.IgnoreFields(Controller{}, "Logger")))
			}
		})
	}
}

func TestController_Create(t *testing.T) {
	tests := []struct {
		name    string
		Logger  *logrus.Logger
		Storage StorageProvider
		model   Bet
		wantErr error
	}{
		{
			name:    "expect success given valid bet",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			model: Bet{
				ID:       uuid.New().String(),
				TableID:  uuid.New().String(),
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:   "expect fail given storage error",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			model: Bet{
				ID:       uuid.New().String(),
				TableID:  uuid.New().String(),
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrCreate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage)
			_, err := c.Create(context.Background(), tt.model)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_List(t *testing.T) {
	tests := []struct {
		name     string
		Logger   *logrus.Logger
		Storage  StorageProvider
		wantBets []Bet
		wantErr  error
	}{
		{
			name:   "expect success given no Bets found",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenList: []Bet{},
			},
			wantBets: []Bet{},
			wantErr:  nil,
		},
		{
			name:   "expect success given Bets found",
			Logger: logrus.New(),

			Storage: mockStorage{
				GivenList: []Bet{
					{
						ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
						TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
						Bet:      "10",
						Type:     TypeStraight,
						Amount:   100,
						Currency: "GBP",
					},
				},
			},
			wantBets: []Bet{
				{
					ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
					TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
					Bet:      "10",
					Type:     TypeStraight,
					Amount:   100,
					Currency: "GBP",
				},
			},
			wantErr: nil,
		},
		{
			name:   "expect fail given storage error",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
				GivenList:  []Bet{},
			},
			wantBets: []Bet{},
			wantErr:  ErrList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage)
			bets, err := c.List(context.Background(), uuid.New().String())

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(bets, tt.wantBets) {
				t.Error(cmp.Diff(bets, tt.wantBets))
			}
		})
	}
}

func TestController_Update(t *testing.T) {
	tests := []struct {
		name    string
		Logger  *logrus.Logger
		Storage StorageProvider
		model   Bet
		wantErr error
	}{
		{
			name:    "expect success given valid bet",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			model: Bet{
				ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:   "expect fail given storage error",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			model: Bet{
				ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrUpdate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage)
			_, err := c.Update(context.Background(), tt.model)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Delete(t *testing.T) {
	tests := []struct {
		name    string
		Logger  *logrus.Logger
		Storage StorageProvider
		id      string
		wantErr error
	}{
		{
			name:    "expect success given valid bet",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			id:      uuid.New().String(),
			wantErr: nil,
		},
		{
			name:   "expect fail given storage error",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			id:      uuid.New().String(),
			wantErr: ErrDelete,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage)
			_, err := c.Delete(context.Background(), uuid.New().String(), tt.id)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Play(t *testing.T) {
	tests := []struct {
		name    string
		Logger  *logrus.Logger
		Storage StorageProvider
		tableID string
		want    Result
		wantErr error
	}{
		{
			name:   "expect success given valid input",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenList: []Bet{},
			},
			tableID: "8117bb87-148c-4fb1-8971-a2d4373b3f19",
			want: Result{
				Winners: []Winner{},
			},
			wantErr: nil,
		},
		{
			name:   "expect success given valid input with found bets",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenList: []Bet{
					{
						ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
						TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
						Bet:      "10",
						Type:     "straight",
						Amount:   1000,
						Currency: "GBP",
					},
				},
			},
			tableID: "8117bb87-148c-4fb1-8971-a2d4373b3f19",
			want: Result{
				Winners: []Winner{
					{
						BetID:    "8117bb87-148c-4fb1-8971-a2d4373b3f19",
						Amount:   1000,
						Currency: "GBP",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "expect fail given storage error",
			Logger: logrus.New(),
			Storage: mockStorage{
				GivenList:  []Bet{},
				GivenError: errors.New("foo"),
			},
			tableID: "8117bb87-148c-4fb1-8971-a2d4373b3f19",
			want:    Result{},
			wantErr: ErrList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage)

			got, err := c.Play(context.Background(), tt.tableID)
			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(Result{}, "Number", "Color")) {
				t.Error(cmp.Diff(err, tt.want, cmpopts.IgnoreFields(Result{}, "Number", "Color")))
			}
		})
	}
}

func TestController_getColor(t *testing.T) {
	tests := []struct {
		name    string
		Logger  *logrus.Logger
		Storage StorageProvider
		number  int
		want    string
	}{
		{
			name:    "expect green given 0",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			number:  0,
			want:    "green",
		},
		{
			name:    "expect black given 2",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			number:  2,
			want:    "black",
		},
		{
			name:    "expect red given 3",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			number:  3,
			want:    "red",
		},
		{
			name:    "expect black given 13",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			number:  13,
			want:    "black",
		},
		{
			name:    "expect red given 14",
			Logger:  logrus.New(),
			Storage: mockStorage{},
			number:  14,
			want:    "red",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage)

			got := c.getColor(tt.number)

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

type mockStorage struct {
	GivenList  []Bet
	GivenID    string
	GivenError error
}

func (m mockStorage) Delete(_ context.Context, _, _ string) (string, error) {
	return m.GivenID, m.GivenError
}

func (m mockStorage) Update(_ context.Context, _ Bet) (string, error) {
	return m.GivenID, m.GivenError
}

func (m mockStorage) List(_ context.Context, _ string, _ ...Bet) ([]Bet, error) {
	return m.GivenList, m.GivenError
}

func (m mockStorage) Create(_ context.Context, _ Bet) (string, error) {
	return m.GivenID, m.GivenError
}

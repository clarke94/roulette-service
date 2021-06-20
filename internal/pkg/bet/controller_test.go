package bet

import (
	"context"
	"errors"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/go-playground/validator/v10"
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
		model     Bet
		wantErr   error
	}{
		{
			name:      "expect success given valid bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.New(),
				TableID:  uuid.New(),
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect fail given no bet table ID",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given no money",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:      uuid.New(),
				TableID: uuid.New(),
				Bet:     "10",
				Type:    TypeStraight,
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
			model: Bet{
				ID:       uuid.New(),
				TableID:  uuid.New(),
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
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		wantBets  []Bet
		wantErr   error
	}{
		{
			name:      "expect success given no Bets found",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList: []Bet{},
			},
			wantBets: []Bet{},
			wantErr:  nil,
		},
		{
			name:      "expect success given Bets found",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList: []Bet{
					{
						ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
						TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
						Bet:      "10",
						Type:     TypeStraight,
						Amount:   100,
						Currency: "GBP",
					},
				},
			},
			wantBets: []Bet{
				{
					ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
					TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
					Bet:      "10",
					Type:     TypeStraight,
					Amount:   100,
					Currency: "GBP",
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
				GivenList:  []Bet{},
			},
			wantBets: []Bet{},
			wantErr:  ErrList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)
			bets, err := c.List(context.Background(), uuid.New())

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(bets, tt.wantBets, cmpopts.IgnoreFields(money.Money{}, "amount", "currency")) {
				t.Error(cmp.Diff(bets, tt.wantBets, cmpopts.IgnoreFields(money.Money{}, "amount", "currency")))
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
		model     Bet
		wantErr   error
	}{
		{
			name:      "expect success given valid bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect fail given no ID",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given no bet ID",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "10",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given no money",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:      uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID: uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:     "10",
				Type:    TypeStraight,
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
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
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
			c := New(tt.Logger, tt.Storage, tt.Validator)
			_, err := c.Update(context.Background(), tt.model)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Delete(t *testing.T) {
	tests := []struct {
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		id        uuid.UUID
		wantErr   error
	}{
		{
			name:      "expect success given valid bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			id:        uuid.New(),
			wantErr:   nil,
		},
		{
			name:      "expect fail given no ID",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			id:        uuid.Nil,
			wantErr:   ErrValidation,
		},
		{
			name:      "expect fail given storage error",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			id:      uuid.New(),
			wantErr: ErrDelete,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)
			_, err := c.Delete(context.Background(), uuid.New(), tt.id)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_CreateUpdateValidate(t *testing.T) {
	tests := []struct {
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		model     Bet
		wantErr   error
	}{
		{
			name:      "expect fail given unknown type",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "10",
				Type:     "foo",
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect fail given red/black type with invalid bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "foo",
				Type:     TypeRedBlack,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: ErrValidation,
		},
		{
			name:      "expect success given red/black type with red bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "red",
				Type:     TypeRedBlack,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect success given red/black type with black bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "black",
				Type:     TypeRedBlack,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect success given straight type with 36 bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "36",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
		{
			name:      "expect success given straight type with number between 0-36 bet",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			model: Bet{
				ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
				Bet:      "12",
				Type:     TypeStraight,
				Amount:   100,
				Currency: "GBP",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)
			_, createErr := c.Create(context.Background(), tt.model)
			_, updateErr := c.Update(context.Background(), tt.model)

			if !cmp.Equal(createErr, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(createErr, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(updateErr, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(updateErr, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Play(t *testing.T) {
	tests := []struct {
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		tableID   uuid.UUID
		want      Result
		wantErr   error
	}{
		{
			name:      "expect success given valid input",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList: []Bet{},
			},
			tableID: uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
			want: Result{
				Winners: []Winner{},
			},
			wantErr: nil,
		},
		{
			name:      "expect success given valid input with found bets",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList: []Bet{
					{
						ID:       uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
						TableID:  uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
						Bet:      "10",
						Type:     "straight",
						Amount:   1000,
						Currency: "GBP",
					},
				},
			},
			tableID: uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
			want: Result{
				Winners: []Winner{
					{
						BetID:    uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
						Amount:   1000,
						Currency: "GBP",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:      "expect fail given storage error",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage: mockStorage{
				GivenList:  []Bet{},
				GivenError: errors.New("foo"),
			},
			tableID: uuid.Must(uuid.Parse("8117bb87-148c-4fb1-8971-a2d4373b3f19")),
			want:    Result{},
			wantErr: ErrList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)

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
		name      string
		Logger    *logrus.Logger
		Validator *validator.Validate
		Storage   StorageProvider
		number    int
		want      string
	}{
		{
			name:      "expect green given 0",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			number:    0,
			want:      "green",
		},
		{
			name:      "expect black given 2",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			number:    2,
			want:      "black",
		},
		{
			name:      "expect red given 3",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			number:    3,
			want:      "red",
		},
		{
			name:      "expect black given 13",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			number:    13,
			want:      "black",
		},
		{
			name:      "expect red given 14",
			Logger:    logrus.New(),
			Validator: validator.New(),
			Storage:   mockStorage{},
			number:    14,
			want:      "red",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.Logger, tt.Storage, tt.Validator)

			got := c.getColor(tt.number)

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

type mockStorage struct {
	GivenList  []Bet
	GivenID    uuid.UUID
	GivenError error
}

func (m mockStorage) Delete(_ context.Context, _, _ uuid.UUID) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

func (m mockStorage) Update(_ context.Context, _ Bet) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

func (m mockStorage) List(_ context.Context, _ uuid.UUID, _ ...Bet) ([]Bet, error) {
	return m.GivenList, m.GivenError
}

func (m mockStorage) Create(_ context.Context, _ Bet) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}

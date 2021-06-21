package test

import (
	"context"
	"github.com/clarke94/roulette-service/internal/pkg/bet"
	storage "github.com/clarke94/roulette-service/storage/bet"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
)

func TestBetStorage_Create(t *testing.T) {
	tests := []struct {
		name    string
		model   bet.Bet
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name: "expect success given valid bet",
			model: bet.Bet{
				ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				Bet:      "foo",
				Type:     "bar",
				Amount:   10,
				Currency: "GBP",
			},
			ctx:     context.Background(),
			want:    "8117bb87-148c-4fb1-8971-a2d4373b3f19",
			wantErr: false,
		},
		{
			name: "expect fail given id already exists",
			model: bet.Bet{
				ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			},
			ctx:     nil,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New(db)

			got, err := s.Create(tt.ctx, tt.model)
			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBetStorage_List(t *testing.T) {
	tests := []struct {
		name    string
		tableID string
		filters []bet.Bet
		want    []bet.Bet
		wantErr bool
	}{
		{
			name:    "expect array of bets given valid request",
			tableID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			filters: nil,
			want: []bet.Bet{
				{
					ID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					TableID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Bet:      "foo",
					Type:     "bar",
					Amount:   10,
					Currency: "GBP",
				},
			},
			wantErr: false,
		},
		{
			name:    "expect array of bets given valid request with filter",
			tableID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			filters: []bet.Bet{
				{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				},
			},
			want: []bet.Bet{
				{
					ID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					TableID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Bet:      "foo",
					Type:     "bar",
					Amount:   10,
					Currency: "GBP",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New(db)

			got, err := s.List(context.Background(), tt.tableID, tt.filters...)

			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBetStorage_Update(t *testing.T) {
	tests := []struct {
		name    string
		model   bet.Bet
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name: "expect fail given bet doesnt exist",
			model: bet.Bet{
				ID:       uuid.New().String(),
				TableID:  uuid.New().String(),
				Bet:      "foo",
				Type:     "bar",
				Amount:   10,
				Currency: "GBP",
			},
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
		{
			name: "expect success bet exists",
			model: bet.Bet{
				ID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				TableID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Bet:      "foo",
				Type:     "bar",
				Amount:   10,
				Currency: "GBP",
			},
			ctx:     context.Background(),
			want:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New(db)

			got, err := s.Update(tt.ctx, tt.model)
			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBetStorage_Delete(t *testing.T) {
	tests := []struct {
		name    string
		tableID string
		betID   string
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name:    "expect fail given bet doesnt exist",
			tableID: uuid.New().String(),
			betID:   uuid.New().String(),
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "expect success bet exists",
			tableID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			betID:   "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			ctx:     context.Background(),
			want:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New(db)

			got, err := s.Delete(tt.ctx, tt.tableID, tt.betID)
			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

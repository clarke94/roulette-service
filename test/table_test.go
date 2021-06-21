package test

import (
	"context"
	"github.com/clarke94/roulette-service/internal/pkg/table"
	storage "github.com/clarke94/roulette-service/storage/table"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
)

func TestTableStorage_Create(t *testing.T) {
	tests := []struct {
		name    string
		model   table.Table
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name: "expect success given valid table",
			model: table.Table{
				ID:         "cccccccc-cccc-cccc-cccc-cccccccccccc",
				Name:       "",
				MaximumBet: 0,
				MinimumBet: 0,
				Currency:   "GBP",
			},
			ctx:     context.Background(),
			want:    "cccccccc-cccc-cccc-cccc-cccccccccccc",
			wantErr: false,
		},
		{
			name: "expect fail given id already exists",
			model: table.Table{
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

func TestTableStorage_List(t *testing.T) {
	tests := []struct {
		name    string
		want    []table.Table
		wantErr bool
	}{
		{
			name: "expect array of tables given valid request",
			want: []table.Table{
				{
					ID:         "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Name:       "",
					MaximumBet: 0,
					MinimumBet: 0,
					Currency:   "GBP",
				},
				{
					ID:         "cccccccc-cccc-cccc-cccc-cccccccccccc",
					Name:       "",
					MaximumBet: 0,
					MinimumBet: 0,
					Currency:   "GBP",
				},
			},
			wantErr: false,
		},
		{
			name: "expect array of tables given valid request with filter",
			want: []table.Table{
				{
					ID:         "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Name:       "",
					MaximumBet: 0,
					MinimumBet: 0,
					Currency:   "GBP",
				},
				{
					ID:         "cccccccc-cccc-cccc-cccc-cccccccccccc",
					Name:       "",
					MaximumBet: 0,
					MinimumBet: 0,
					Currency:   "GBP",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New(db)

			got, err := s.List(context.Background())

			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestTableStorage_Update(t *testing.T) {
	tests := []struct {
		name    string
		model   table.Table
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name: "expect fail given table doesnt exist",
			model: table.Table{
				ID:         uuid.New().String(),
				Name:       "",
				MaximumBet: 0,
				MinimumBet: 0,
				Currency:   "GBP",
			},
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
		{
			name: "expect success table exists",
			model: table.Table{
				ID:         "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				Name:       "",
				MaximumBet: 0,
				MinimumBet: 0,
				Currency:   "GBP",
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

func TestTableStorage_Delete(t *testing.T) {
	tests := []struct {
		name    string
		tableID string
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name:    "expect fail given table doesnt exist",
			tableID: uuid.New().String(),
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "expect success table exists",
			tableID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			ctx:     context.Background(),
			want:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New(db)

			got, err := s.Delete(tt.ctx, tt.tableID)
			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

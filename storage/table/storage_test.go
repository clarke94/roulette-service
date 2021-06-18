package table

import (
	"context"
	"github.com/clarke94/roulette-service/internal/pkg/table"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want Storage
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Create(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx   context.Context
		model table.Table
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				DB: tt.fields.DB,
			}
			got, err := s.Create(tt.args.ctx, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

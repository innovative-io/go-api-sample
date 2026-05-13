package services

import (
	"context"
	"flag"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/innovative-io/go-api-sample/cmd/tests"
	"github.com/innovative-io/go-api-sample/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkDogInserts(b *testing.B) {
	teardownTests := tests.SetupTests(b, postgres.Open(tests.ConnectionString))
	defer teardownTests(b)

	dogsService := NewDogsService(tests.DB)

	for i := 0; i < b.N; i++ {
		dog := &models.Dog{
			ID:        uuid.New(),
			Name:      tests.RandString(12),
			Breed:     tests.RandString(12),
			Color:     tests.RandString(12),
			Birthdate: time.Now(),
			Weight:    rand.Intn(298) + 1,
		}
		_, err := dogsService.Add(context.Background(), dog)
		if err != nil {
			panic(err)
		}
	}
}

func TestIntegration_NewDogsService(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want DogsService
	}{
		{
			name: "Should get valid interface back",
			args: args{db: tests.DB},
			want: &dogsService{db: tests.DB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDogsService(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDogsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dogsService_Add(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		dog *models.Dog
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dogsService{
				db: tt.fields.db,
			}
			got, err := s.Add(tt.args.ctx, tt.args.dog)
			if (err != nil) != tt.wantErr {
				t.Errorf("dogsService.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dogsService.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegration_dogsService_Add(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	dogID := uuid.New()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		dog *models.Dog
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			name:   "Should add a dog to the database",
			fields: fields{db: tests.DB},
			args: args{
				ctx: context.Background(),
				dog: &models.Dog{
					ID:        dogID,
					Name:      "Snowball",
					Breed:     "Sniba Inu",
					Color:     "Cream",
					Birthdate: time.Now(),
					Weight:    22,
				},
			},
			want:    &dogID,
			wantErr: false,
		},
		{
			name:   "Should not be able to add empty dog",
			fields: fields{db: tests.DB},
			args: args{
				ctx: context.Background(),
				dog: &models.Dog{
					ID: uuid.New(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dogsService{
				db: tt.fields.db,
			}
			got, err := s.Add(tt.args.ctx, tt.args.dog)
			if (err != nil) != tt.wantErr {
				t.Errorf("dogsService.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dogsService.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegration_dogsService_Delete(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Should delete dog by ID",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: tests.Dogs[0].ID},
			wantErr: false,
		},
		{
			name:    "Should not delete an ID that does not exists",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: uuid.New()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dogsService{
				db: tt.fields.db,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("dogsService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntegration_dogsService_Get(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Should get all the dogs",
			fields:    fields{db: tests.DB},
			args:      args{ctx: context.Background()},
			wantCount: 4,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dogsService{
				db: tt.fields.db,
			}
			got, err := s.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("dogsService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("dogsService.Get() = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}
func TestIntegration_dogsService_GetOne(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Dog
		wantErr bool
	}{
		{
			name:    "Should get dog by ID",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: tests.Dogs[0].ID},
			want:    &tests.Dogs[0],
			wantErr: false,
		},
		{
			name:    "Should not get dog that does not exist",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: uuid.New()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dogsService{
				db: tt.fields.db,
			}
			got, err := s.GetOne(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("dogsService.GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got.ID, tt.want.ID) {
					t.Errorf("dogsService.GetOne() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestIntegration_dogsService_Update(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
		dog *models.Dog
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Should update a dog by ID",
			fields: fields{db: tests.DB},
			args: args{ctx: context.Background(), id: tests.Dogs[0].ID, dog: &models.Dog{
				Name:      "Snowball",
				Breed:     "Shiba Inu",
				Color:     "Cream",
				Birthdate: time.Date(2020, 6, 12, 0, 0, 0, 0, time.UTC),
				Weight:    22,
			}},
			wantErr: false,
		},
		{
			name:   "Should not update a dog with no valid ID",
			fields: fields{db: tests.DB},
			args: args{ctx: context.Background(), id: uuid.New(), dog: &models.Dog{
				Name:      "Snowball",
				Breed:     "Shiba Inu",
				Color:     "Cream",
				Birthdate: time.Date(2020, 6, 12, 0, 0, 0, 0, time.UTC),
				Weight:    22,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dogsService{
				db: tt.fields.db,
			}
			if err := s.Update(tt.args.ctx, tt.args.id, tt.args.dog); (err != nil) != tt.wantErr {
				t.Errorf("dogsService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

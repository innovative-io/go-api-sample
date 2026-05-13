package services

import (
	"context"
	"flag"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/innovative-io/go-api-sample/cmd/tests"
	"github.com/innovative-io/go-api-sample/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkCatInserts(b *testing.B) {
	teardownTests := tests.SetupTests(b, postgres.Open(tests.ConnectionString))
	defer teardownTests(b)

	catsService := NewCatsService(tests.DB)

	for i := 0; i < b.N; i++ {
		cat := &models.Cat{
			ID:        uuid.New(),
			Name:      tests.RandString(12),
			Breed:     tests.RandString(12),
			Color:     tests.RandString(12),
			Birthdate: time.Now(),
			Weight:    rand.Intn(98) + 1,
		}
		_, err := catsService.Add(context.Background(), cat)
		if err != nil {
			panic(err)
		}
	}
}

func TestIntegration_NewCatsService(t *testing.T) {
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
		want CatsService
	}{
		{
			name: "Should get valid interface back",
			args: args{db: tests.DB},
			want: &catsService{db: tests.DB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCatsService(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCatsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_catsService_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{Conn: db},
	})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	zeroID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	birth, _ := time.Parse(time.RFC3339, "2020-02-10T00:00:00Z")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		cat *models.Cat
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			name:   "Add cat to database",
			fields: fields{db: gdb},
			args: args{ctx: context.Background(), cat: &models.Cat{
				Name:      "Nacho",
				Breed:     "Tabby",
				Color:     "Orange",
				Birthdate: birth,
				Weight:    17,
			}},
			want:    &zeroID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "cats"`).WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			s := &catsService{
				db: tt.fields.db,
			}
			got, err := s.Add(tt.args.ctx, tt.args.cat)
			if (err != nil) != tt.wantErr {
				t.Errorf("catsService.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("catsService.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegration_catsService_Add(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	catID := uuid.New()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		cat *models.Cat
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			name:   "Should add a cat to the database",
			fields: fields{db: tests.DB},
			args: args{
				ctx: context.Background(),
				cat: &models.Cat{
					ID:        catID,
					Name:      "Nacho",
					Breed:     "Tabby",
					Color:     "Orange",
					Birthdate: time.Now(),
					Weight:    17,
				},
			},
			want:    &catID,
			wantErr: false,
		},
		{
			name:   "Should not be able to add empty cat",
			fields: fields{db: tests.DB},
			args: args{
				ctx: context.Background(),
				cat: &models.Cat{
					ID: uuid.New(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &catsService{
				db: tt.fields.db,
			}
			got, err := s.Add(tt.args.ctx, tt.args.cat)
			if (err != nil) != tt.wantErr {
				t.Errorf("catsService.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("catsService.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_catsService_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{Conn: db},
	})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

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
			name:    "Delete cat from database",
			fields:  fields{db: gdb},
			args:    args{ctx: context.Background(), id: uuid.New()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`DELETE FROM "cats" WHERE`).WithArgs(tt.args.id).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			s := &catsService{
				db: tt.fields.db,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("catsService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntegration_catsService_Delete(t *testing.T) {
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
			name:    "Should delete cat by ID",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: tests.Cats[0].ID},
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
			s := &catsService{
				db: tt.fields.db,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("catsService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_catsService_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{Conn: db},
	})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Cat
		wantErr bool
	}{
		{
			name:    "Get all cats",
			fields:  fields{db: gdb},
			args:    args{ctx: context.Background()},
			want:    []models.Cat{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{""}))
			s := &catsService{
				db: tt.fields.db,
			}
			got, err := s.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("catsService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("catsService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegration_catsService_Get(t *testing.T) {
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
			name:      "Should get all the cats",
			fields:    fields{db: tests.DB},
			args:      args{ctx: context.Background()},
			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &catsService{
				db: tt.fields.db,
			}
			got, err := s.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("catsService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("catsService.Get() = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}
func TestIntegration_catsService_GetOne(t *testing.T) {
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
		want    *models.Cat
		wantErr bool
	}{
		{
			name:    "Should get cat by ID",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: tests.Cats[0].ID},
			want:    &tests.Cats[0],
			wantErr: false,
		},
		{
			name:    "Should not get cat that does not exist",
			fields:  fields{db: tests.DB},
			args:    args{ctx: context.Background(), id: uuid.New()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &catsService{
				db: tt.fields.db,
			}
			got, err := s.GetOne(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("catsService.GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got.ID, tt.want.ID) {
					t.Errorf("catsService.GetOne() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestIntegration_catsService_Update(t *testing.T) {
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
		cat *models.Cat
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Should update a cat by ID",
			fields: fields{db: tests.DB},
			args: args{ctx: context.Background(), id: tests.Cats[0].ID, cat: &models.Cat{
				Name:      "Nacho",
				Breed:     "Tabby",
				Color:     "Orange",
				Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
				Weight:    19,
			}},
			wantErr: false,
		},
		{
			name:   "Should not update a cat with no valid ID",
			fields: fields{db: tests.DB},
			args: args{ctx: context.Background(), id: uuid.New(), cat: &models.Cat{
				Name:      "Nacho",
				Breed:     "Tabby",
				Color:     "Orange",
				Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
				Weight:    19,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &catsService{
				db: tt.fields.db,
			}
			if err := s.Update(tt.args.ctx, tt.args.id, tt.args.cat); (err != nil) != tt.wantErr {
				t.Errorf("catsService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_catsService_Count(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{Conn: db},
	})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "Count all cats",
			fields:  fields{db: gdb},
			args:    args{ctx: context.Background()},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count\(\*\)`).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
			s := &catsService{
				db: tt.fields.db,
			}
			got, err := s.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("catsService.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("catsService.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

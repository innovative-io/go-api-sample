package controllers

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
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

	router := NewRouter(tests.DB)

	for i := 0; i < b.N; i++ {
		cat := &models.Cat{
			ID:        uuid.New(),
			Name:      tests.RandString(12),
			Breed:     tests.RandString(12),
			Color:     tests.RandString(12),
			Birthdate: time.Now(),
			Weight:    rand.Intn(98) + 1,
		}

		data, _ := json.Marshal(cat)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/cats", bytes.NewReader(data))
		req.Header.Add("Content-type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			panic("failed to create cat")
		}
	}
}

func TestCatsDelete(t *testing.T) {
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

	router := NewRouter(gdb)

	testID := uuid.New()

	type args struct {
		method   string
		endpoint string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name: "Should delete a cat by ID",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/cats/%s", testID.String()),
			},
			wantResponse: fmt.Sprintf("{\"deleted\":\"%s\"}", testID.String()),
			wantCode:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`DELETE FROM "cats" WHERE`).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
			router.ServeHTTP(w, req)

			if tt.wantCode != w.Code {
				t.Errorf("CatsDelete() error = %v, wantCode %v", w.Code, tt.wantCode)
				return
			}

			if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
				t.Errorf("CatsDelete() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
			}
		})
	}
}

func TestIntegrationCatsDelete(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name: "Should delete a cat by ID",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/cats/%s", tests.Cats[0].ID.String()),
			},
			wantResponse: fmt.Sprintf("{\"deleted\":\"%s\"}", tests.Cats[0].ID.String()),
			wantCode:     http.StatusOK,
		},
		{
			name: "Should not delete an invalid ID",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/cats/%s", "not_a_valid_id"),
			},
			wantResponse: "{\"message\":\"invalid id\"}",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not delete an ID that does not exist",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/cats/%s", uuid.New().String()),
			},
			wantResponse: "{\"message\":\"not found\"}",
			wantCode:     http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("CatsDelete() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
			t.Errorf("CatsDelete() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

func TestIntegrationCatsGet(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantCode  int
	}{
		{
			name: "Should all cats",
			args: args{
				method:   "GET",
				endpoint: "/cats",
			},
			wantCount: len(tests.Cats),
			wantCode:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("CatsGet() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		cats := make([]models.Cat, 0)
		err := json.Unmarshal(w.Body.Bytes(), &cats)
		if err != nil {
			t.Errorf("CatsGet() error = %v, wantCount %v", err, tt.wantCount)
			return
		}

		if tt.wantCount != len(cats) {
			t.Errorf("CatsGet() error = %v, wantCount %v", len(cats), tt.wantCount)
		}
	}
}

func TestIntegrationCatsGetOne(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse *models.Cat
		wantCode     int
	}{
		{
			name: "Should get a cat by ID",
			args: args{
				method:   "GET",
				endpoint: fmt.Sprintf("/cats/%s", tests.Cats[0].ID.String()),
			},
			wantResponse: &tests.Cats[0],
			wantCode:     http.StatusOK,
		},
		{
			name: "Should not get an invalid ID",
			args: args{
				method:   "GET",
				endpoint: fmt.Sprintf("/cats/%s", "not_a_valid_id"),
			},
			wantResponse: nil,
			wantCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("CatsGetOne() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if tt.wantResponse != nil {
			cat := new(models.Cat)
			err := json.Unmarshal(w.Body.Bytes(), &cat)
			if err != nil {
				t.Errorf("CatsGetOne() error = %v, wantCount %v", err, tt.wantResponse)
				return
			}

			if !reflect.DeepEqual(tt.wantResponse.ID, cat.ID) {
				t.Errorf("CatsGetOne() error = %v, wantCode %v", cat, tt.wantResponse)
			}
		}
	}
}

func TestIntegrationCatsPost(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
		body     *models.Cat
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name: "Should add a cat to the database",
			args: args{
				method:   "POST",
				endpoint: "/cats",
				body: &models.Cat{
					Name:      "Fluffy",
					Breed:     "Bengal",
					Color:     "Orange",
					Birthdate: time.Date(2022, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    5,
				},
			},
			wantResponse: "id",
			wantCode:     http.StatusCreated,
		},
		{
			name: "Should not add a cat with empty breed to the database",
			args: args{
				method:   "POST",
				endpoint: "/cats",
				body: &models.Cat{
					Name:      "Fluffy",
					Breed:     "",
					Color:     "Orange",
					Birthdate: time.Date(2022, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    5,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not add an invalid cat to the database",
			args: args{
				method:   "POST",
				endpoint: "/cats",
				body: &models.Cat{
					Name:      "Fluffy",
					Breed:     "Bengal",
					Color:     "Orange",
					Birthdate: time.Date(2022, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    200,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.args.body)
		if err != nil {
			t.Errorf("CatsPost() unable to marshal %v", tt.args.body)
			return
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, bytes.NewReader(data))
		req.Header.Add("Content-type", "application/json")
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("CatsPost() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !strings.Contains(w.Body.String(), tt.wantResponse) {
			t.Errorf("CatsPost() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

func TestIntegrationCatsPut(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
		body     *models.Cat
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name: "Should update a cat",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/cats/%s", tests.Cats[0].ID.String()),
				body: &models.Cat{
					Name:      "Nacho",
					Breed:     "Tabby",
					Color:     "Orange",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    20,
				},
			},
			wantResponse: "id",
			wantCode:     http.StatusAccepted,
		},
		{
			name: "Should not update an invalid id",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/cats/%s", "invalid_id_here"),
				body: &models.Cat{
					Name:      "Nacho",
					Breed:     "Tabby",
					Color:     "Orange",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    20,
				},
			},
			wantResponse: "invalid",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not update a cat with empty color",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/cats/%s", tests.Cats[0].ID.String()),
				body: &models.Cat{
					Name:      "Nacho",
					Breed:     "Tabby",
					Color:     "",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    20,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not update a cat with empty color",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/cats/%s", tests.Cats[0].ID.String()),
				body: &models.Cat{
					Name:      "Nacho",
					Breed:     "Tabby",
					Color:     "Orange",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    200,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.args.body)
		if err != nil {
			t.Errorf("CatsPut() unable to marshal %v", tt.args.body)
			return
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, bytes.NewReader(data))
		req.Header.Add("Content-type", "application/json")
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("CatsPut() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !strings.Contains(w.Body.String(), tt.wantResponse) {
			t.Errorf("CatsPut() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

func TestCatsCount(t *testing.T) {
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

	router := NewRouter(gdb)

	tests := []struct {
		name         string
		wantResponse string
		wantCode     int
	}{
		{
			name:         "Should return count of cats",
			wantResponse: `{"count":5}`,
			wantCode:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count\(\*\)`).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/cats/count", nil)
			router.ServeHTTP(w, req)

			if tt.wantCode != w.Code {
				t.Errorf("CatsCount() error = %v, wantCode %v", w.Code, tt.wantCode)
				return
			}

			if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
				t.Errorf("CatsCount() error = %v, want %v", w.Body.String(), tt.wantResponse)
			}
		})
	}
}

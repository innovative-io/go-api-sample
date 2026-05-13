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

func BenchmarkDogInserts(b *testing.B) {
	teardownTests := tests.SetupTests(b, postgres.Open(tests.ConnectionString))
	defer teardownTests(b)

	router := NewRouter(tests.DB)

	for i := 0; i < b.N; i++ {
		dog := &models.Dog{
			ID:        uuid.New(),
			Name:      tests.RandString(12),
			Breed:     tests.RandString(12),
			Color:     tests.RandString(12),
			Birthdate: time.Now(),
			Weight:    rand.Intn(98) + 1,
		}

		data, _ := json.Marshal(dog)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/dogs", bytes.NewReader(data))
		req.Header.Add("Content-type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			panic("failed to create dog")
		}
	}
}

func TestDogsDelete(t *testing.T) {
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
			name: "Should delete a dog by ID",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/dogs/%s", testID.String()),
			},
			wantResponse: fmt.Sprintf("{\"deleted\":\"%s\"}", testID.String()),
			wantCode:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`DELETE FROM "dogs" WHERE`).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
			router.ServeHTTP(w, req)

			if tt.wantCode != w.Code {
				t.Errorf("DogsDelete() error = %v, wantCode %v", w.Code, tt.wantCode)
				return
			}

			if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
				t.Errorf("DogsDelete() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
			}
		})
	}
}

func TestIntegrationDogsDelete(t *testing.T) {
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
			name: "Should delete a dog by ID",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/dogs/%s", tests.Dogs[0].ID.String()),
			},
			wantResponse: fmt.Sprintf("{\"deleted\":\"%s\"}", tests.Dogs[0].ID.String()),
			wantCode:     http.StatusOK,
		},
		{
			name: "Should not delete an invalid ID",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/dogs/%s", "not_a_valid_id"),
			},
			wantResponse: "{\"message\":\"invalid id\"}",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not delete an ID that does not exist",
			args: args{
				method:   "DELETE",
				endpoint: fmt.Sprintf("/dogs/%s", uuid.New().String()),
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
			t.Errorf("DogsDelete() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
			t.Errorf("DogsDelete() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

func TestIntegrationDogsGet(t *testing.T) {
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
			name: "Should all dogs",
			args: args{
				method:   "GET",
				endpoint: "/dogs",
			},
			wantCount: len(tests.Dogs),
			wantCode:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("DogsGet() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		dogs := make([]models.Dog, 0)
		err := json.Unmarshal(w.Body.Bytes(), &dogs)
		if err != nil {
			t.Errorf("DogsGet() error = %v, wantCount %v", err, tt.wantCount)
			return
		}

		if tt.wantCount != len(dogs) {
			t.Errorf("DogsGet() error = %v, wantCount %v", len(dogs), tt.wantCount)
		}
	}
}

func TestIntegrationDogsGetOne(t *testing.T) {
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
		wantResponse *models.Dog
		wantCode     int
	}{
		{
			name: "Should get a dog by ID",
			args: args{
				method:   "GET",
				endpoint: fmt.Sprintf("/dogs/%s", tests.Dogs[0].ID.String()),
			},
			wantResponse: &tests.Dogs[0],
			wantCode:     http.StatusOK,
		},
		{
			name: "Should not get an invalid ID",
			args: args{
				method:   "GET",
				endpoint: fmt.Sprintf("/dogs/%s", "not_a_valid_id"),
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
			t.Errorf("DogsGetOne() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if tt.wantResponse != nil {
			dog := new(models.Dog)
			err := json.Unmarshal(w.Body.Bytes(), &dog)
			if err != nil {
				t.Errorf("DogsGetOne() error = %v, wantCount %v", err, tt.wantResponse)
				return
			}

			if !reflect.DeepEqual(tt.wantResponse.ID, dog.ID) {
				t.Errorf("DogsGetOne() error = %v, wantCode %v", dog, tt.wantResponse)
			}
		}
	}
}

func TestIntegrationDogsPost(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
		body     *models.Dog
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name: "Should add a dog to the database",
			args: args{
				method:   "POST",
				endpoint: "/dogs",
				body: &models.Dog{
					Name:      "Spike",
					Breed:     "Bulldog",
					Color:     "Grey",
					Birthdate: time.Date(2021, 12, 10, 0, 0, 0, 0, time.UTC),
					Weight:    55,
				},
			},
			wantResponse: "id",
			wantCode:     http.StatusCreated,
		},
		{
			name: "Should not add a dog with empty breed to the database",
			args: args{
				method:   "POST",
				endpoint: "/dogs",
				body: &models.Dog{
					Name:      "Spike",
					Breed:     "",
					Color:     "Grey",
					Birthdate: time.Date(2021, 12, 10, 0, 0, 0, 0, time.UTC),
					Weight:    55,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not add a dog with invalid weight to the database",
			args: args{
				method:   "POST",
				endpoint: "/dogs",
				body: &models.Dog{
					Name:      "Spike",
					Breed:     "Bulldog",
					Color:     "Grey",
					Birthdate: time.Date(2021, 12, 10, 0, 0, 0, 0, time.UTC),
					Weight:    0,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.args.body)
		if err != nil {
			t.Errorf("DogsPost() unable to marshal %v", tt.args.body)
			return
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, bytes.NewReader(data))
		req.Header.Add("Content-type", "application/json")
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("DogsPost() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !strings.Contains(w.Body.String(), tt.wantResponse) {
			t.Errorf("DogsPost() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

func TestIntegrationDogsPut(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
		body     *models.Dog
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name: "Should update a dog",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/dogs/%s", tests.Dogs[0].ID.String()),
				body: &models.Dog{
					Name:      "0111",
					Breed:     "Pitbull Nix",
					Color:     "White/Brindle",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    65,
				},
			},
			wantResponse: "id",
			wantCode:     http.StatusAccepted,
		},
		{
			name: "Should not update an invalid id",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/dogs/%s", "invalid_id_here"),
				body: &models.Dog{
					Name:      "0111",
					Breed:     "Pitbull Nix",
					Color:     "White/Brindle",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    65,
				},
			},
			wantResponse: "invalid",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not update a dog with empty color",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/dogs/%s", tests.Dogs[0].ID.String()),
				body: &models.Dog{
					Name:      "0111",
					Breed:     "Pitbull Nix",
					Color:     "",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    65,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "Should not update a dog with invalid weight",
			args: args{
				method:   "PUT",
				endpoint: fmt.Sprintf("/dogs/%s", tests.Dogs[0].ID.String()),
				body: &models.Dog{
					Name:      "0111",
					Breed:     "Pitbull Nix",
					Color:     "White/Brindle",
					Birthdate: time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
					Weight:    301,
				},
			},
			wantResponse: "",
			wantCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.args.body)
		if err != nil {
			t.Errorf("DogsPut() unable to marshal %v", tt.args.body)
			return
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, bytes.NewReader(data))
		req.Header.Add("Content-type", "application/json")
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("DogsPut() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !strings.Contains(w.Body.String(), tt.wantResponse) {
			t.Errorf("DogsPut() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

func TestDogsCount(t *testing.T) {
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
			name:         "Should return count of dogs",
			wantResponse: `{"count":7}`,
			wantCode:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count\(\*\)`).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/dogs/count", nil)
			router.ServeHTTP(w, req)

			if tt.wantCode != w.Code {
				t.Errorf("DogsCount() error = %v, wantCode %v", w.Code, tt.wantCode)
				return
			}

			if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
				t.Errorf("DogsCount() error = %v, want %v", w.Body.String(), tt.wantResponse)
			}
		})
	}
}

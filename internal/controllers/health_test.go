package controllers

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/innovative-io/go-api-sample/cmd/tests"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHealthGet(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error opening stub db: %s", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.Dialector{Config: &postgres.Config{Conn: db}})
	if err != nil {
		t.Fatalf("unexpected error opening gorm db: %s", err)
	}

	router := NewRouter(gdb)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HealthGet() status = %d, want %d", w.Code, http.StatusOK)
	}
	if want := `{"message":"ok"}`; w.Body.String() != want {
		t.Errorf("HealthGet() body = %q, want %q", w.Body.String(), want)
	}
}

func TestIntegrationHealthGet(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	teardownTests := tests.SetupTests(t, postgres.Open(tests.ConnectionString))
	defer teardownTests(t)

	router := NewRouter(tests.DB)

	type args struct {
		method   string
		endpoint string
		body     interface{}
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
		wantCode     int
	}{
		{
			name:         "Should pass health check",
			args:         args{method: "GET", endpoint: "/health", body: nil},
			wantResponse: "{\"message\":\"ok\"}",
			wantCode:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
		router.ServeHTTP(w, req)

		if tt.wantCode != w.Code {
			t.Errorf("HealthGet() error = %v, wantCode %v", w.Code, tt.wantCode)
			return
		}

		if !reflect.DeepEqual(tt.wantResponse, w.Body.String()) {
			t.Errorf("HealthGet() error = %v, wantCode %v", w.Body.String(), tt.wantResponse)
		}
	}
}

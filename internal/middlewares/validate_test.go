package middlewares

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestIntegrationValidateHeader(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"ok"}`))
	})
	handler := ValidateHeader(mux)

	type args struct {
		method   string
		endpoint string
		header   map[string]string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Should check health",
			args: args{
				method:   "GET",
				endpoint: "/health",
				header:   map[string]string{},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Should not be able check health",
			args: args{
				method:   "GET",
				endpoint: "/health",
				header: map[string]string{
					"X-Not-Valid": "true",
				},
			},
			wantCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
		for key, value := range tt.args.header {
			req.Header.Add(key, value)
		}
		handler.ServeHTTP(w, req)
		if w.Code != tt.wantCode {
			t.Errorf("%s: expected status %d, got %d", tt.name, tt.wantCode, w.Code)
		}
	}
}

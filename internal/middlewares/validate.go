package middlewares

import (
	"encoding/json"
	"net/http"
)

func ValidateHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Not-Valid") != "" {
			data, _ := json.Marshal(map[string]any{"message": "not allowed"})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write(data)
			return
		}
		next.ServeHTTP(w, r)
	})
}

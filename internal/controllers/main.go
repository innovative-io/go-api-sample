package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/innovative-io/go-api-sample/docs"
	"github.com/innovative-io/go-api-sample/internal/middlewares"
	"github.com/innovative-io/go-api-sample/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

type Router struct {
	cats     services.CatsService
	dogs     services.DogsService
	validate *validator.Validate
	db       *gorm.DB
}

func NewRouter(db *gorm.DB) http.Handler {
	r := &Router{
		cats:     services.NewCatsService(db),
		dogs:     services.NewDogsService(db),
		validate: validator.New(),
		db:       db,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", r.HealthGet)

	mux.HandleFunc("DELETE /cats/{id}", r.CatsDelete)
	mux.HandleFunc("GET /cats/count", r.CatsCount)
	mux.HandleFunc("GET /cats", r.CatsGet)
	mux.HandleFunc("GET /cats/{id}", r.CatsGetOne)
	mux.HandleFunc("POST /cats", r.CatsPost)
	mux.HandleFunc("PUT /cats/{id}", r.CatsPut)

	mux.HandleFunc("DELETE /dogs/{id}", r.DogsDelete)
	mux.HandleFunc("GET /dogs/count", r.DogsCount)
	mux.HandleFunc("GET /dogs", r.DogsGet)
	mux.HandleFunc("GET /dogs/{id}", r.DogsGetOne)
	mux.HandleFunc("POST /dogs", r.DogsPost)
	mux.HandleFunc("PUT /dogs/{id}", r.DogsPut)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return middlewares.LoggingMiddleware(corsMiddleware(middlewares.ValidateHeader(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (r *Router) bind(w http.ResponseWriter, req *http.Request, v any) error {
	req.Body = http.MaxBytesReader(w, req.Body, 1<<20)
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}
	return r.validate.Struct(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/innovative-io/go-api-sample/internal/controllers"
	"github.com/innovative-io/go-api-sample/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var version string = "development"

// @title Go API Sample
// @version 1.0
// @description This is a sample API in go
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	printVersion()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := setupDatabase(getConnectionString())
	router := controllers.NewRouter(db)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func getConnectionString() string {
	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		panic("unable to get a connection string")
	}
	return connectionString
}

func printVersion() {
	fmt.Printf("Starting go-api-sample %s\n\n", version)
}

func setupDatabase(connectionString string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err))
	}

	if err := db.AutoMigrate(&models.Cat{}, &models.Dog{}); err != nil {
		panic(err)
	}

	return db
}

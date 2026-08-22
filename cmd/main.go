package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/http/handlers"
	"github.com/tynou/avito-assignment/internal/http/middleware"
	"github.com/tynou/avito-assignment/internal/repo/users"
	"github.com/tynou/avito-assignment/internal/service/auth"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		log.Fatalf("failed to initialize migration: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration error: %v", err)
	}
	defer m.Close()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer pool.Close()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	validate := validator.New()

	userRepo := users.NewUserRepo(pool)

	authService := auth.NewAuthService(userRepo)

	apiHandler := handlers.NewApiHandler(validate, authService)

	router := http.NewServeMux()

	router.Handle("GET /api/info", middleware.AuthMiddleware(http.HandlerFunc(apiHandler.Info)))
	router.Handle("POST /api/sendCoin", middleware.AuthMiddleware(http.HandlerFunc(apiHandler.SendCoin)))
	router.Handle("GET /api/buy/{item}", middleware.AuthMiddleware(http.HandlerFunc(apiHandler.Buy)))

	router.HandleFunc("POST /api/auth", apiHandler.Auth)

	log.Printf("server running on port: %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

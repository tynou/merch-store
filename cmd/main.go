package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/config"
	"github.com/tynou/avito-assignment/internal/http/handlers"
	"github.com/tynou/avito-assignment/internal/http/middleware"
	"github.com/tynou/avito-assignment/internal/repo/merch"
	"github.com/tynou/avito-assignment/internal/repo/purchases"
	"github.com/tynou/avito-assignment/internal/repo/users"
	"github.com/tynou/avito-assignment/internal/service/auth"
	"github.com/tynou/avito-assignment/internal/service/purchase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
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

	validate := validator.New()

	userRepo := users.NewUserRepo(pool)
	merchRepo := merch.NewMerchRepo(pool)
	purchaseRepo := purchases.NewPurchaseRepo(pool)

	authService := auth.NewAuthService(cfg, userRepo)
	purchaseService := purchase.NewPurchaseService(merchRepo, purchaseRepo)

	apiHandler := handlers.NewApiHandler(validate, authService, purchaseService)

	router := http.NewServeMux()

	router.Handle("GET /api/info", middleware.AuthMiddleware(cfg, http.HandlerFunc(apiHandler.Info)))
	router.Handle("POST /api/sendCoin", middleware.AuthMiddleware(cfg, http.HandlerFunc(apiHandler.SendCoin)))
	router.Handle("GET /api/buy/{item}", middleware.AuthMiddleware(cfg, http.HandlerFunc(apiHandler.Buy)))

	router.HandleFunc("POST /api/auth", apiHandler.Auth)

	log.Printf("server running on port: %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

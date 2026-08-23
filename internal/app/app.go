package app

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/tynou/avito-assignment/internal/repo/transfers"
	"github.com/tynou/avito-assignment/internal/repo/users"
	"github.com/tynou/avito-assignment/internal/service/auth"
	"github.com/tynou/avito-assignment/internal/service/info"
	"github.com/tynou/avito-assignment/internal/service/purchase"
	"github.com/tynou/avito-assignment/internal/service/transfer"
)

func Initialize(ctx context.Context, cfg *config.Config, dsn string, migrationPath string) (http.Handler, func(), error) {
	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize migration: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, nil, fmt.Errorf("migration error: %w", err)
	}
	m.Close()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("db connection error: %w", err)
	}

	validate := validator.New()

	userRepo := users.NewUserRepo(pool)
	merchRepo := merch.NewMerchRepo(pool)
	purchaseRepo := purchases.NewPurchaseRepo(pool)
	transferRepo := transfers.NewTransferRepo(pool)

	authService := auth.NewAuthService(cfg, userRepo)
	purchaseService := purchase.NewPurchaseService(merchRepo, purchaseRepo)
	infoService := info.NewInfoService(userRepo, purchaseRepo, transferRepo)
	transferService := transfer.NewTransferService(userRepo, transferRepo)

	apiHandler := handlers.NewApiHandler(validate, authService, purchaseService, infoService, transferService)

	router := http.NewServeMux()

	router.Handle("GET /api/info", middleware.AuthMiddleware(cfg, http.HandlerFunc(apiHandler.Info)))
	router.Handle("POST /api/sendCoin", middleware.AuthMiddleware(cfg, http.HandlerFunc(apiHandler.SendCoin)))
	router.Handle("GET /api/buy/{item}", middleware.AuthMiddleware(cfg, http.HandlerFunc(apiHandler.Buy)))
	router.HandleFunc("POST /api/auth", apiHandler.Auth)

	return router, pool.Close, nil
}

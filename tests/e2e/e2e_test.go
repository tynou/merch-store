package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/tynou/avito-assignment/internal/app"
	"github.com/tynou/avito-assignment/internal/config"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/http/handlers"
)

var (
	testServer *httptest.Server
	testClient *http.Client
	queries    *db.Queries
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		postgres.WithDatabase("test_shop"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	dsn, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("test db connection error: %v", err)
	}

	queries = db.New(pool)

	mockCfg := &config.Config{
		Port:   "8080",
		JWTKey: []byte("my_super_secret_key"),
	}

	router, cleanup, err := app.Initialize(ctx, mockCfg, dsn, "file://../../db/migrations")
	if err != nil {
		log.Fatalf("failed to initialize app for tests: %v", err)
	}

	testServer = httptest.NewServer(router)
	testClient = testServer.Client()

	code := m.Run()

	testServer.Close()
	pool.Close()
	cleanup()
	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}

func cleanDB(t *testing.T) {
	t.Helper()

	err := queries.CleanTestData(context.Background())

	require.NoError(t, err)
}

func authenticate(t *testing.T, username, password string) string {
	t.Helper()

	authRequest, err := json.Marshal(handlers.AuthRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("failed to marshal auth request: %v", err)
	}

	resp, err := http.Post(testServer.URL+"/api/auth", "application/json", bytes.NewBuffer(authRequest))
	if err != nil {
		t.Fatalf("failed to send auth request: %v", err)
	}
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResponse handlers.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("failed to decode auth response: %v", err)
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, authResponse.Token)

	return authResponse.Token
}

func buyMerch(t *testing.T, token, item string) {
	t.Helper()

	req, err := http.NewRequest("GET", testServer.URL+"/api/buy/"+item, nil)
	if err != nil {
		t.Fatalf("failed to create buy merch request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send buy merch request: %v", err)
	}
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func sendCoin(t *testing.T, token, toUser string, amount int32) {
	t.Helper()

	sendCoinRequest, err := json.Marshal(handlers.SendCoinRequest{
		ToUser: toUser,
		Amount: amount,
	})
	if err != nil {
		t.Fatalf("failed to marshal send coin request: %v", err)
	}

	req, err := http.NewRequest("POST", testServer.URL+"/api/sendCoin", bytes.NewBuffer(sendCoinRequest))
	if err != nil {
		t.Fatalf("failed to create send coin request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send send coin request: %v", err)
	}
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func getInfo(t *testing.T, token string) handlers.InfoResponse {
	t.Helper()

	req, err := http.NewRequest("GET", testServer.URL+"/api/info", nil)
	if err != nil {
		t.Fatalf("failed to create get info request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send get info request: %v", err)
	}
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var infoResponse handlers.InfoResponse
	err = json.NewDecoder(resp.Body).Decode(&infoResponse)
	if err != nil {
		t.Fatalf("failed to decode info response: %v", err)
	}

	assert.NoError(t, err)

	return infoResponse
}

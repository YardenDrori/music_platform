package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/YardenDrori/music-platform/internal/auth"
	"github.com/YardenDrori/music-platform/internal/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//==========Setup==========
	//nolint
	godotenv.Load(".env")

	m, err := migrate.New("file://migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("setting up migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer db.Close()
	slog.Info("connected to database")

	signingKey := os.Getenv("SIGNING_KEY")
	accessTokenDurStr := os.Getenv("ACCESS_TOKEN_DURATION")
	refreshTokenDurStr := os.Getenv("REFRESH_TOKEN_DURATION")
	if signingKey == "" {
		return fmt.Errorf("SIGNING_KEY not present in env")
	}
	if accessTokenDurStr == "" {
		return fmt.Errorf("ACCESS_TOKEN_DURATION not present in env")
	}
	slog.Info("access token duration", "duration", accessTokenDurStr)
	if refreshTokenDurStr == "" {
		return fmt.Errorf("REFRESH_TOKEN_DURATION not present in env")
	}
	slog.Info("refresh token duration", "duration", refreshTokenDurStr)

	accessTokenDur, err := time.ParseDuration(accessTokenDurStr)
	if err != nil {
		return fmt.Errorf("invalid ACCESS_TOKEN_DURATION: %w", err)
	}
	slog.Info("access token duration", "duration", accessTokenDur)
	refreshTokenDur, err := time.ParseDuration(refreshTokenDurStr)
	if err != nil {
		return fmt.Errorf("invalid REFRESH_TOKEN_DURATION: %w", err)
	}
	slog.Info("refresh token duration", "duration", refreshTokenDur)

	//==========User==========
	userHandler := user.NewHandler(
		user.NewService(
			user.NewPostgresRepository(db),
		),
	)

	//==========Server==========
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.Login)

	log.Println("server starting on :8080")
	return http.ListenAndServe(":8080", mux)
}

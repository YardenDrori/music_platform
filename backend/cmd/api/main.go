package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/YardenDrori/music-platform/internal/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//==========Setup==========
	godotenv.Load("../../.env")

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer db.Close()

	signingKey := os.Getenv("SIGNING_KEY")
	accessTokenDurStr := os.Getenv("ACCESS_TOKEN_DURATION")
	refreshTokenDurStr := os.Getenv("REFRESH_TOKEN_DURATION")
	if signingKey == "" {
		return fmt.Errorf("SIGNING_KEY not present in env")
	}
	if accessTokenDurStr == "" {
		return fmt.Errorf("ACCESS_TOKEN_DURATION not present in env")
	}
	if refreshTokenDurStr == "" {
		return fmt.Errorf("REFRESH_TOKEN_DURATION not present in env")
	}

	accessTokenDur, err := time.ParseDuration(accessTokenDurStr)
	if err != nil {
		return fmt.Errorf("invalid ACCESS_TOKEN_DURATION: %w", err)
	}
	refreshTokenDur, err := time.ParseDuration(refreshTokenDurStr)
	if err != nil {
		return fmt.Errorf("invalid REFRESH_TOKEN_DURATION: %w", err)
	}

	//==========User==========
	userHandler := user.NewHandler(
		user.NewService(
			user.NewPostgresRepository(db),
			user.NewJwtTokenizer([]byte(signingKey), accessTokenDur, refreshTokenDur),
		),
	)

	//==========Server==========
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	log.Println("server starting on :8080")
	return http.ListenAndServe(":8080", mux)
}

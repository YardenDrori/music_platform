package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/YardenDrori/music-platform/internal/user"
)

func main() {
	//==========Setup==========
	godotenv.Load("../../.env")

	//==========UserApis==========
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	defer db.Close()

	userRepo := user.NewPostgresRepository(db)

	signingKey := os.Getenv("SIGNING_KEY")
	accessTokenDurStr := os.Getenv("ACCESS_TOKEN_DURATION")
	refreshTokenDurStr := os.Getenv("REFRESH_TOKEN_DURATION")
	if signingKey == "" {
		panic("SIGNING_KEY not present in env")
	}
	if accessTokenDurStr == "" {
		panic("ACCESS_TOKEN_DURATION not present in env")
	}
	if refreshTokenDurStr == "" {
		panic("REFRESH_TOKEN_DURATION not present in env")
	}

	accessTokenDur, err := time.ParseDuration(accessTokenDurStr)
	if err != nil {
		panic("invalid access token duration syntax in env")
	}
	refreshTokenDur, err := time.ParseDuration(refreshTokenDurStr)
	if err != nil {
		panic("invalid refresh token duration syntax in env")
	}

	userTokenizer := user.NewJwtTokenizer(
		[]byte(signingKey),
		accessTokenDur,
		refreshTokenDur,
	)

	UserApis := user.NewHandler(
		user.NewService(
			userRepo,
			userTokenizer,
		),
	)

	//==========INIT==========
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", UserApis.Register)
	mux.HandleFunc("POST /login", UserApis.Login)

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

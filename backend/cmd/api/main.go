package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/tools/go/analysis/passes/stringintconv"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/auth"
	"github.com/YardenDrori/music-platform/internal/storage"
	"github.com/YardenDrori/music-platform/internal/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type router struct {
	mux        *http.ServeMux
	middleware []func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request, error)
}

func routeWithMiddleware(
	mux *http.ServeMux,
	middleware ...func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request, error),
) *router {
	router := router{
		mux:        mux,
		middleware: nil,
	}
	for _, v := range middleware {
		router.middleware = append(router.middleware, v)
	}
	return &router
}

func (router *router) route(
	pattern string,
	handler func(http.ResponseWriter, *http.Request) error,
) {
	router.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		var err error
		for _, v := range router.middleware {
			w, r, err = v(w, r)
			if err != nil {
				apperrors.HandlerError(w, r, err)
				return
			}
		}
		if err := handler(w, r); err != nil {
			apperrors.HandlerError(w, r, err)
		}
	})
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

	storageEndpoint := os.Getenv("STORAGE_ENDPOINT")
	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	storageSecretAccessKey := os.Getenv("STORAGE_SECRET_ACCESS_KEY")
	storageIsSecureStr := os.Getenv("STORAGE_IS_SECURE")
	if storageEndpoint == "" {
		return fmt.Errorf("storage endpoint not present in env: %w", err)
	}
	if storageAccessKey == "" {
		return fmt.Errorf("storage access key not present in env: %w", err)
	}
	if storageSecretAccessKey == "" {
		return fmt.Errorf("storage secret access key not present in env: %w", err)
	}
	if storageIsSecureStr == "" {
		return fmt.Errorf("storage is secure not present in env")
	}
	storageIsSecure, err := strconv.ParseBool(storageIsSecureStr)
	if err != nil {
		return fmt.Errorf("STORAGE_IS_SECURE value format invalid")
	}

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
	userService := user.NewService(user.NewPostgresRepository(db), user.NewArgon2idPasswordHasher())
	userHandler := user.NewHandler(userService)

	//==========Auth==========
	authService := auth.NewService(
		auth.NewPostgresRepository(db),
		auth.NewJwtTokenizer([]byte(signingKey), accessTokenDur, refreshTokenDur),
		userService,
	)
	authHandler := auth.NewHandler(authService)
	requireAuth := auth.NewRequireAuth(authService)

	//==========STORAGE==========
	storageCreds := credentials.NewStaticV4(storageAccessKey, storageSecretAccessKey, "")
	storageOpts := &minio.Options{
		Creds:  storageCreds,
		Secure: storageIsSecure,
	}
	minioClient, err := minio.New(storageEndpoint, storageOpts)
	if err != nil {
		return fmt.Errorf("failed to instantiate minioClient: %w", err)
	}
	sorageService := storage.NewService(minioClient)

	//==========ROUTER==========
	mux := http.NewServeMux()
	root := routeWithMiddleware(mux)
	authn := routeWithMiddleware(mux, requireAuth)

	//==========Server==========
	root.route("POST /api/register", authHandler.Register)
	root.route("POST /api/login", authHandler.Login)
	root.route("POST /api/token", authHandler.GetAccessToken)
	authn.route("GET /api/me", userHandler.GetMe)
	mux.HandleFunc("PATCH /api/me", requireAuth(userHandler.UpdateMe))
	mux.HandleFunc("DELETE /api/me", requireAuth(userHandler.DisableMe))

	log.Println("server starting on :8080")
	return http.ListenAndServe(":8080", mux)
}

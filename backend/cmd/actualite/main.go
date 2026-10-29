package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"actualiteco/actualite/internal/core/logger"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"actualiteco/actualite/internal/core/transport/server"
	auth_repository "actualiteco/actualite/internal/features/auth/repository/postgres"
	auth_service "actualiteco/actualite/internal/features/auth/service"
	auth_transport "actualiteco/actualite/internal/features/auth/transport/http"
)

// @title actualiteco/actualite
// @version 1.0
// @host 127.0.0.1:8080
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error with reading .env")
		os.Exit(1)
	}

	log, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("error with initializing new logger")
		os.Exit(1)
	}
	defer log.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pool, err := postgres_pool.NewConnectionPool(ctx, postgres_pool.NewConfigMust())
	if err != nil {
		log.Fatal("creating new pgxPool", zap.Error(err))
	}
	defer pool.Close()

	serv := server.NewHTTPServer(server.NewConfigMust(), log)
	serv.RegisterSwagger()

	log.Debug("initializing auth service")
	authRepo := auth_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepo)
	authHandler := auth_transport.NewAuthHandler(authService)

	serv.RegisterRoutes(authHandler.Routes()...)

	if err = serv.Run(ctx); err != nil {
		log.Error("server error occured", zap.Error(err))
	}
}

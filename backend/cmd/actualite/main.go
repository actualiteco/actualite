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
	auth_service "actualiteco/actualite/internal/features/auth/service"
	auth_transport "actualiteco/actualite/internal/features/auth/transport/http"
	post_repository "actualiteco/actualite/internal/features/post/repository/postgres"
	post_service "actualiteco/actualite/internal/features/post/service"
	post_transport "actualiteco/actualite/internal/features/post/transport/http"
	users_repository "actualiteco/actualite/internal/features/profile/repository/postgres"
	profile_service "actualiteco/actualite/internal/features/profile/service"
	profile_transport "actualiteco/actualite/internal/features/profile/transport/http"
)

// @title actualiteco/actualite
// @version 1.0
// @host 127.0.0.1:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer" и JWT-токен через пробел.
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
	usersRepo := users_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(usersRepo)
	authHandler := auth_transport.NewAuthHandler(authService)

	log.Debug("initializing profile service")
	profileService := profile_service.NewProfileService(usersRepo)
	profileHandler := profile_transport.NewProfileHTTPHandler(profileService)

	log.Debug("initializing post service")
	postRepository := post_repository.New(pool)
	postService := post_service.New(postRepository)
	postHandler := post_transport.NewPostsHandler(postService)

	serv.RegisterRoutes(authHandler, profileHandler, postHandler)

	if err = serv.Run(ctx); err != nil {
		log.Error("server error occured", zap.Error(err))
	}
}

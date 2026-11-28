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
	"actualiteco/actualite/internal/core/middleware"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	smtpmailhogclient "actualiteco/actualite/internal/core/smtp/mailhog"
	"actualiteco/actualite/internal/core/storage"
	tokenservice "actualiteco/actualite/internal/core/token/service"
	"actualiteco/actualite/internal/core/transport/server"
	sessionsrepo "actualiteco/actualite/internal/features/auth/repository/postgres"
	auth_service "actualiteco/actualite/internal/features/auth/service"
	auth_transport "actualiteco/actualite/internal/features/auth/transport/http"
	notification_repository "actualiteco/actualite/internal/features/notification/repository/postgres"
	post_repository "actualiteco/actualite/internal/features/post/repository/postgres"
	post_service "actualiteco/actualite/internal/features/post/service"
	post_transport "actualiteco/actualite/internal/features/post/transport/http"
	users_repository "actualiteco/actualite/internal/features/profile/repository/postgres"
	profile_service "actualiteco/actualite/internal/features/profile/service"
	profile_transport "actualiteco/actualite/internal/features/profile/transport/http"
	search_repository "actualiteco/actualite/internal/features/search/repository/postgres"
	search_service "actualiteco/actualite/internal/features/search/service"
	search_transport "actualiteco/actualite/internal/features/search/transport/http"
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

	log.Debug("initializing tx manager")
	txManager := postgres_pool.NewTxManager(pool)

	log.Debug("initializing token service")
	tokenCfg, err := tokenservice.NewConfig()
	if err != nil {
		log.Fatal("creating token config", zap.Error(err))
	}
	tokenService := tokenservice.New(tokenCfg)

	authMiddleware := middleware.AuthMiddleware(tokenService)
	serverCfg := server.NewConfigMust()
	serv := server.NewHTTPServer(serverCfg, log, authMiddleware)
	serv.RegisterSwagger()

	log.Debug("initializing smtp client")
	smtpClient := smtpmailhogclient.NewMust()

	log.Debug("initializing auth service")
	usersRepo := users_repository.NewAuthRepository(pool)
	sessionsRepo := sessionsrepo.New(pool)
	authService, err := auth_service.NewAuthService(usersRepo, sessionsRepo, tokenService, smtpClient, txManager, serverCfg.FrontendBaseURL)
	if err != nil {
		log.Fatal("init auth service (service layer)", zap.Error(err))
	}
	authHandler := auth_transport.NewAuthHandler(authService)

	log.Debug("initializing file storage")
	storageClient := storage.NewSeaweedFSClient(storage.NewConfigMust())

	log.Debug("initializing profile service")
	notificationRepo := notification_repository.New(pool)
	profileService := profile_service.NewProfileService(usersRepo, notificationRepo, storageClient, txManager)
	profileHandler := profile_transport.NewProfileHTTPHandler(profileService)

	log.Debug("initializing post service")
	postRepository := post_repository.New(pool)
	postService := post_service.New(postRepository, txManager, storageClient)
	postHandler := post_transport.NewPostsHandler(postService)

	log.Debug("initializing search service")
	searchRepository := search_repository.New(pool)
	searchService := search_service.New(searchRepository, txManager)
	searchHandler := search_transport.NewSearchHTTPHandler(searchService)

	serv.RegisterRoutes(authHandler, profileHandler, postHandler, searchHandler)

	if err = serv.Run(ctx); err != nil {
		log.Error("server error occured", zap.Error(err))
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/tryingmyb3st/PolyTweet/internal/core/logger"
	postgres_pool "github.com/tryingmyb3st/PolyTweet/internal/core/repository/postgres"
	"github.com/tryingmyb3st/PolyTweet/internal/core/transport/server"
)

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

	if err = serv.Run(ctx); err != nil {
		log.Error("server error occured", zap.Error(err))
	}
}

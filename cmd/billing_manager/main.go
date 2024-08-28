package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mongo_billing_repository "github.com/ThePositree/billing_manager/internal/adapter/repository/billing/mongo"
	mongo_user_repository "github.com/ThePositree/billing_manager/internal/adapter/repository/user/mongo"
	"github.com/ThePositree/billing_manager/internal/config"
	http_controller "github.com/ThePositree/billing_manager/internal/controller/http"
	"github.com/ThePositree/billing_manager/internal/usecase/billing_managing/billing_managing_std"
	"github.com/ThePositree/billing_manager/internal/usecase/user_managing/user_managing_std"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGKILL,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create config")
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create mongo connect")
	}

	userRepo, err := mongo_user_repository.New(ctx, logger, mongoClient, mongo_user_repository.Config{
		Database:   cfg.Database,
		Collection: cfg.UserCollection,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create user repo")
	}

	billingRepo, err := mongo_billing_repository.New(ctx, logger, mongoClient, mongo_billing_repository.Config{
		Database:   cfg.Database,
		Collection: cfg.BillingCollection,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create billing repo")
	}

	billingManaging := billing_managing_std.New(userRepo, billingRepo)
	userManaging := user_managing_std.New(userRepo)

	ctrl := http_controller.New(logger, billingManaging, userManaging, cfg.HttpPort, cfg.AdminPassword)

	go func() {
		<-ctx.Done()
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error().Err(err).Msg("Mongo disconnect")
		}
	}()

	logger.Info().Msg(fmt.Sprintf("HTTP controller started on %d port", cfg.HttpPort))
	ctrl.Start(ctx)
}

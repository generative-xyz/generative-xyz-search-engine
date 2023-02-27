package main

import (
	"context"
	projectMongoDb "generative-xyz-search-engine/internal/repository/mongo"
	usecase "generative-xyz-search-engine/internal/usecase/indexer"
	"generative-xyz-search-engine/pkg/core"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/driver/mongodb"
	"generative-xyz-search-engine/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	mongoDb *mongo.Database
	err     error
)

func main() {
	app := core.NewWorkerApp(
		"search-engine-worker",
		onInit,
		onClose,
	)
	defer onClose(app)
	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func onInit(app *core.App) {
	// init tracer
	tracer.Start(
		tracer.WithEnv(viper.GetString("ENV")),
		tracer.WithService(viper.GetString("APP_NAME")),
		tracer.WithLogger(logger.AtLog),
	)

	algoliaClient := algolia.NewAlgoliaClient()
	// init database
	mongoDb, err = mongodb.Init()
	if err != nil {
		logger.AtLog.Fatalf("connect mongodb failed: %v", err)
	}

	projectRepo := projectMongoDb.NewProjectRepository(mongoDb)
	tokenUriRepo := projectMongoDb.NewTokenRepository(mongoDb)

	ch := make(chan struct{}, 1)
	projectUc := usecase.NewProjectIndexerUsecase(algoliaClient, projectRepo, tokenUriRepo, ch)

	go func() {
		projectUc.Schedule()
	}()
}

func onClose(app *core.App) {
	_ = mongoDb.Client().Disconnect(context.Background())
	tracer.Stop()
	logger.AtLog.Info("search-engine-worker STOPPED")
}

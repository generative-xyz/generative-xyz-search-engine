package main

import (
	"context"
	repoMongoDb "generative-xyz-search-engine/internal/repository/mongo"
	usecase "generative-xyz-search-engine/internal/usecase/indexer"
	"generative-xyz-search-engine/pkg/core"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/driver/mongodb"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/pkg/redis"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	redisDb redis.Client
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

	// inti redis
	redisDb = redis.NewClient()

	projectRepo := repoMongoDb.NewProjectRepository(mongoDb)
	tokenUriRepo := repoMongoDb.NewTokenRepository(mongoDb)
	userRepo := repoMongoDb.NewUserRepository(mongoDb)
	dexBtcListingRepo := repoMongoDb.NewDexBtcListingRepository(mongoDb)

	ch := make(chan struct{}, 1)
	projectUc := usecase.NewProjectIndexerUsecase(algoliaClient, redisDb, projectRepo, tokenUriRepo, userRepo, dexBtcListingRepo, ch)

	go func() {
		projectUc.Schedule()
	}()

	go func() {
		<-ch
		logger.AtLog.Error("Dead goroutine detected. Graceful termination in 30 sec.")
		time.Sleep(30 * time.Second)
		panic("Panic worker because dead goroutine detected.")
	}()
}

func onClose(app *core.App) {
	_ = mongoDb.Client().Disconnect(context.Background())
	tracer.Stop()
	logger.AtLog.Info("search-engine-worker STOPPED")
}

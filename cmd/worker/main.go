package main

import (
	"generative-xyz-search-engine/pkg/core"
	"generative-xyz-search-engine/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

}

func onClose(app *core.App) {
	tracer.Stop()
	logger.AtLog.Info("search-engine-worker STOPPED")
}

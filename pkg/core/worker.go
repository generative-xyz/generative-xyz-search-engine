package core

import (
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/utils"

	"github.com/spf13/viper"
)

// NewWorkerApp returns new app
func NewWorkerApp(
	name string,
	onInit OnInitCallback,
	onClose OnCloseCallback) *App {
	app := App{
		Name: name,
		Type: AppTypeWorker,
	}
	// init default flags, then parse
	app.InitDefaultFlags()
	app.ParseFlags()
	app.InitCommon()
	app.Version = viper.GetString("APP_VERSION")

	logger.AtLog.Infof("%v. Starting...", app.Info())

	// handle sigterm
	utils.HandleSigterm(func() {
		logger.AtLog.Infof("%v. Stopping...", app.Info())
		onClose(&app)
	})

	onInit(&app)

	return &app
}

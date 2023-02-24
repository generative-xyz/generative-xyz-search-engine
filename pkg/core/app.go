package core

import (
	"context"
	"flag"
	"fmt"
	"generative-xyz-search-engine/pkg/conf"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/utils"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

// AppType -- micro app type
type AppType string

// OnCloseCallback --
type OnCloseCallback func(*App)

// OnInitCallback --
type OnInitCallback func(*App)

const (
	// MaxProcs --
	MaxProcs = 3000

	// AppTypeWorker -- worker app
	AppTypeWorker AppType = "worker"

	// DefaultVersion --
	DefaultVersion = "0.1"

	// ConfigSourceFile --
	ConfigSourceFile = "file"
	// ConfigsSourceEnv--
	ConfigSourceEnv = "env"
	// ConfigSourceNone --
	ConfigSourceNone = "none"

	// LoggerFormatJSON --
	LoggerFormatJSON = "json"
	// LoggerFormatText --
	LoggerFormatText = "text"

	// LoggerOutputConsole --
	LoggerOutputConsole = "console"
	// LoggerOutputFile --
	LoggerOutputFile = "file"
)

var (
	// SupportedConfigSources -- list config source supported
	SupportedConfigSources = []string{ConfigSourceFile, ConfigSourceEnv, ConfigSourceNone}

	// SupportedLoggerFormats -- list logger format supported
	SupportedLoggerFormats = []string{LoggerFormatJSON, LoggerFormatText}

	// SupportedLoggerOutputs -- list logger output supported
	SupportedLoggerOutputs = []string{LoggerOutputConsole}
)

// App -- micro app
type App struct {
	// Common
	Name        string
	DisplayName string
	Version     string
	Type        AppType
	ctx         context.Context
	cancel      context.CancelFunc
	// Config
	ConfigLoaded bool
	ConfigSource string
	ConfigFile   string
	// Logger
	LoggerFormat      string
	LoggerOutput      string
	LoggerEnableDebug bool
}

// InitCommon -- initializes common things: config, logger, port, services discovery, tracing
func (app *App) InitCommon() {
	app.PreCheck()
	app.InitConfig()
	app.InitLogger()
	app.setMaxGR()
}

// InitDefaultFlags -- init default flags of app
func (app *App) InitDefaultFlags() {
	// flags for config
	flag.StringVar(&app.ConfigSource, "config-source", "env", "Configuration source: file or env or none")
	flag.StringVar(&app.ConfigFile, "config-file", "", "Configuration file path")

	// flags for logger
	flag.StringVar(&app.LoggerOutput, "logger-output", "console", "Logger output: console or file")
	flag.StringVar(&app.LoggerFormat, "logger-format", "text", "Logger format: text or json")
	flag.BoolVar(&app.LoggerEnableDebug, "logger-enable-debug", false, "Logger enable debug level")
}

// ParseFlags -- parse flags
func (app *App) ParseFlags() {
	flag.Parse()
}

// Info -- show info of app
func (app *App) Info() string {
	return fmt.Sprintf("App Info. Type: %v. Name: %v. Version %v",
		app.Type, app.Name, app.Version)
}

// PreCheck -- pre check some conditional
func (app *App) PreCheck() {
	// validate app name: check if empty, check valid regex pattern
	if utils.IsStringEmpty(app.Name) {
		log.Panic("must set app name")
	}
	matched, err := regexp.MatchString("^[a-z0-9-]+$", app.Name)
	if err != nil {
		log.Panic(err)
	}
	if !matched {
		log.Panic("app name not match with regex pattern [a-z0-9-]")
	}

	// validate app type
	if utils.IsStringEmpty(string(app.Type)) {
		log.Panic("must set app type")
	}
}

// InitConfig -- initializes config
func (app *App) InitConfig() {
	configSource := utils.StringTrimSpace(strings.ToLower(app.ConfigSource))

	// if configSource is empty, set default to `file` type
	if utils.IsStringEmpty(configSource) {
		configSource = ConfigSourceFile
	}

	// check supported configType
	if !slices.Contains(SupportedConfigSources, configSource) {
		log.Panicf("Not support config source: %v", configSource)
	}

	result := false
	switch configSource {
	case ConfigSourceFile:
		configFile := utils.StringTrimSpace(app.ConfigFile)
		// if configFile is empty, get from os env
		if utils.IsStringEmpty(configFile) {
			configFile = conf.GetConfigLocation()
		}
		if utils.IsStringEmpty(configFile) {
			// read by default
			result = conf.ReadConfig("config", ".", "./conf", "./config")
		} else {
			// read by input file
			result = conf.ReadConfigByFile(configFile)
		}
	case ConfigSourceEnv:
		result = conf.ReadConfigFromEnvVariables()
	case ConfigSourceNone:
		log.Println("Not using config file")
		result = true
	}

	if !result {
		log.Panic("Could not load config")
	}

	app.ConfigLoaded = true
	log.Printf("Config loaded from %s \n", configSource)
}

// InitLogger -- initializes logger
func (app *App) InitLogger() {
	// log output
	loggerOutput := conf.GetLogOutput()
	if utils.IsStringEmpty(loggerOutput) {
		loggerOutput = utils.StringTrimSpace(strings.ToLower(app.LoggerOutput))
		if utils.IsStringEmpty(loggerOutput) {
			loggerOutput = LoggerOutputConsole
		}
	}
	// check supported loggerType
	if !slices.Contains(SupportedLoggerOutputs, loggerOutput) {
		log.Panicf("Not supported logger type: %v", loggerOutput)
	}
	app.LoggerOutput = loggerOutput

	// log format
	loggerFormat := conf.GetLogFortmat()
	if utils.IsStringEmpty(loggerFormat) {
		loggerFormat = utils.StringTrimSpace(strings.ToLower(app.LoggerFormat))
		if utils.IsStringEmpty(loggerFormat) {
			loggerFormat = LoggerFormatText
		}
	}
	// check supported loggerFormat
	if !slices.Contains(SupportedLoggerFormats, loggerFormat) {
		log.Panicf("Not supported logger format: %v", loggerFormat)
	}
	app.LoggerFormat = loggerFormat

	// overwrite log enable debug
	if conf.GetLogEnableDebug() {
		app.LoggerEnableDebug = conf.GetLogEnableDebug()
	}

	switch app.LoggerOutput {
	case LoggerOutputConsole:
		switch app.LoggerFormat {
		case LoggerFormatJSON:
			logger.InitLoggerDefault(app.LoggerEnableDebug)
		case LoggerFormatText:
			logger.InitLoggerDefaultDev()
		}
	}
	logger.AtLog.Info("Logger loaded")
}

func (app *App) setMaxGR() {
	max := viper.GetInt("ROOT_MAX_PROCS")
	if max <= 0 {
		max = MaxProcs
	}
	runtime.GOMAXPROCS(max)
}

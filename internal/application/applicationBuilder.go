package application

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/sudzekai/web-os-api/packages/abstractions"
	"github.com/sudzekai/web-os-api/packages/executor"
	"github.com/sudzekai/web-os-api/packages/logging"
	"github.com/sudzekai/web-os-api/packages/modulesloader"
	"github.com/sudzekai/web-os-api/packages/server"
)

type ApplicationBuilder struct {
	loggerFactory abstractions.ILoggerFactory
	configuration abstractions.IConfiguration
}

func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{
		loggerFactory: logging.NewLoggerFactory(nil),
	}
}

func (builder *ApplicationBuilder) SetOutput(w io.Writer) {
	builder.loggerFactory.SetWriter(w)
}

func (builder *ApplicationBuilder) SetConfiguration(cfg abstractions.IConfiguration) {
	builder.configuration = cfg
}

func (builder *ApplicationBuilder) GetLoggerFactory() abstractions.ILoggerFactory {
	return builder.loggerFactory
}

func (builder *ApplicationBuilder) Build() (*application, error) {
	log := builder.loggerFactory.NewLogger("sys:builder")

	cfgBytes, _ := json.MarshalIndent(
		builder.configuration.GetSettings(),
		"",
		"    ",
	)
	cfgStr := string(cfgBytes)

	re := regexp.MustCompile(`"(?i)(user|password)"\s*:\s*"[^"]*"`)
	cfgStr = re.ReplaceAllString(cfgStr, `"${1}": "***"`)

	log.LogDebug("сборка приложения с текущей конфигурацией:\n%s", cfgStr)

	// apploggerfactory
	appLoggerFactory := logging.NewLoggerFactory(nil)
	err := builder.configureAppLoggerFactory(appLoggerFactory)

	if err != nil {
		return nil, err
	}

	// server
	srv := server.NewServer(appLoggerFactory)
	err = builder.configureServer(srv)

	if err != nil {
		return nil, err
	}

	// executor
	executor := executor.NewExecutor(appLoggerFactory)

	// modulesloader
	loader := modulesloader.NewModulesLoader(
		appLoggerFactory,
		executor,
		srv.GetRegistry(),
	)

	return &application{
		loggerFactory: appLoggerFactory,
		server:        srv,
		executor:      executor,
		modulesLoader: loader,
		configuration: builder.configuration,
	}, nil
}

func (builder *ApplicationBuilder) configureServer(srv abstractions.IServer) error {
	host := builder.configuration.GetString("webhost.host")
	srv.SetHost(host)

	port := builder.configuration.GetInt("webhost.port")
	srv.SetPort(port)

	return nil
}

func (builder *ApplicationBuilder) configureAppLoggerFactory(loggerFactory abstractions.ILoggerFactory) error {
	logLevel := builder.configuration.GetString("logging.default.loglevel")
	err := loggerFactory.SetMinLevelStr(logLevel)

	if err != nil {
		return err
	}

	loggerFactory.SetWriter(os.Stdout)

	return nil
}

func errNotFound(optionName string) error {
	return fmt.Errorf("%s не был найден в конфигурации или имеет неверный формат данных", optionName)
}

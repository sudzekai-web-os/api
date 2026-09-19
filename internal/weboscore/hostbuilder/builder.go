package hostbuilder

import (
	"fmt"
	"os"
	"strings"

	cfg "github.com/sudzekai-web-os/api/internal/configuration"
	"github.com/sudzekai-web-os/api/internal/plugins"
	"github.com/sudzekai-web-os/api/internal/weboscore/host"
	"github.com/sudzekai-web-os/api/internal/weboscore/loggingbuilder"
	"github.com/sudzekai-web-os/core"
	"github.com/sudzekai-web-os/executor"
	"github.com/sudzekai-web-os/logging"
	"github.com/sudzekai-web-os/server"
)

type HostBuilder struct {
	loggingBuilder *loggingbuilder.LoggingBuilder
	configuration  core.IConfiguration
	logger         core.ILogger
}

func New() *HostBuilder {
	lf := logging.NewLoggerFactory()
	lf.AddWriter(logging.NewConsoleWriter(os.Stdout))
	lf.SetMinLevel(core.LogLevel_INFORMATION)
	return &HostBuilder{
		loggingBuilder: loggingbuilder.New(),
		logger:         lf.NewLogger("host-builder"),
	}
}

func (hb *HostBuilder) SetConfiguration(cfg core.IConfiguration) *HostBuilder {
	hb.configuration = cfg
	return hb
}

func (hb *HostBuilder) Build() *host.Host {
	hostLoggerFactory := hb.getConfiguredLoggerFactory()

	srv := hb.getConfiguredServer(hostLoggerFactory)
	executor := hb.getConfiguredExecutor(hostLoggerFactory)
	loader := hb.getConfiguredPluginsLoader(hostLoggerFactory, executor, srv.GetRegistry())

	return host.New(
		srv,
		hostLoggerFactory,
		executor,
		loader,
		hb.configuration,
	)
}

func (hb *HostBuilder) getConfiguredLoggerFactory() core.ILoggerFactory {
	loggerFactory := logging.NewLoggerFactory()

	consoleEnabled := hb.configuration.GetBool(cfg.LoggingConsoleEnabled)

	levelStr := hb.configuration.GetString(cfg.LoggingLevel)

	level, err := parseLogLevel(*levelStr)

	if err != nil {
		hb.LogAndExit("%s", err.Error())
	}

	loggerFactory.SetMinLevel(level)

	if consoleEnabled == nil {
		hb.LogConfValueNilAndExit(cfg.LoggingConsoleEnabled)
	}

	if *consoleEnabled {
		loggerFactory.AddWriter(logging.NewConsoleWriter(os.Stdout))
	}

	grpcEnabled := hb.configuration.GetBool(cfg.LoggingGrpcEnabled)

	if grpcEnabled == nil {
		hb.LogConfValueNilAndExit(cfg.LoggingGrpcEnabled)
	}

	if *grpcEnabled {
		grpcEndpoint := hb.configuration.GetString(cfg.LoggingGrpcEndpoint)

		if grpcEndpoint == nil {
			hb.LogConfValueNilAndExit(cfg.LoggingGrpcEndpoint)
		}

		loggerFactory.AddWriter(logging.NewGrpcWriter(*grpcEndpoint))
	}

	fileEnabled := hb.configuration.GetBool(cfg.LoggingFileEnabled)

	if fileEnabled == nil {
		hb.LogConfValueNilAndExit(cfg.LoggingFileEnabled)
	}

	if *fileEnabled {
		fileName := hb.configuration.GetString(cfg.LoggingFileName)

		if fileName == nil {
			hb.LogConfValueNilAndExit(cfg.LoggingFileName)
		}

		loggerFactory.AddWriter(logging.NewFileWriter(*fileName))
	}

	return loggerFactory
}

func (hb *HostBuilder) getConfiguredServer(loggerFactory core.ILoggerFactory) core.IServer {
	srv := server.NewServer(loggerFactory)

	host := hb.configuration.GetString(cfg.WebhostHost)
	if host == nil {
		hb.LogConfValueNilAndExit(cfg.WebhostHost)
	}

	port := hb.configuration.GetInt(cfg.WebhostPort)
	if port == nil {
		hb.LogConfValueNilAndExit(cfg.WebhostPort)
	}

	srv.SetHost(*host)
	srv.SetPort(*port)

	return srv
}

func (hb *HostBuilder) getConfiguredExecutor(loggerFactory core.ILoggerFactory) core.IExecutor {
	executor := executor.NewExecutor(loggerFactory)
	return executor
}

func (hb *HostBuilder) getConfiguredPluginsLoader(
	loggerFactory core.ILoggerFactory,
	executor core.IExecutor,
	registry core.IHandlersRegistry) *plugins.PluginsLoader {
	loader := plugins.NewLoader(
		loggerFactory,
		executor,
		registry,
		hb.configuration,
	)

	return loader
}

func parseLogLevel(level string) (core.LogLevel, error) {
	switch strings.ToLower(level) {
	case "none":
		return core.LogLevel_NONE, nil
	case "debug", "dbug":
		return core.LogLevel_DEBUG, nil
	case "info", "information":
		return core.LogLevel_INFORMATION, nil
	case "warn", "warning":
		return core.LogLevel_WARNING, nil
	case "err", "error":
		return core.LogLevel_ERROR, nil
	case "crit", "critical":
		return core.LogLevel_CRITICAL, nil
	default:
		return core.LogLevel_NONE, fmt.Errorf("неизвестный уровень логирования %s", level)
	}
}

func (hb *HostBuilder) LogAndExit(format string, args ...any) {
	hb.logger.LogCritical(format, args...)
	os.Exit(-1)
}

func (hb *HostBuilder) LogConfValueNilAndExit(propertyName string) {
	hb.logger.LogCritical("configuration property %s value was nil", propertyName)
	os.Exit(-1)
}

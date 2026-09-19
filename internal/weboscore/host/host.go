package host

import (
	"os"

	cfg "github.com/sudzekai-web-os/api/internal/configuration"
	"github.com/sudzekai-web-os/api/internal/middlewares/resultfilter"
	"github.com/sudzekai-web-os/api/internal/plugins"
	"github.com/sudzekai-web-os/core"
)

type Host struct {
	server        core.IServer
	loggerFactory core.ILoggerFactory
	executor      core.IExecutor
	pluginsLoader *plugins.PluginsLoader
	configuration core.IConfiguration
}

func New(
	server core.IServer,
	loggerFactory core.ILoggerFactory,
	executor core.IExecutor,
	pluginsLoader *plugins.PluginsLoader,
	configuration core.IConfiguration,
) *Host {
	return &Host{
		server:        server,
		loggerFactory: loggerFactory,
		executor:      executor,
		pluginsLoader: pluginsLoader,
		configuration: configuration,
	}
}

func (host *Host) Run() {
	log := host.loggerFactory.NewLogger("host")

	log.LogInformation("приложение запущено")

	host.loadPlugins()

	host.server.GetRegistry().SetResultFilter(resultfilter.GetFilterFunc(host.loggerFactory))
	host.server.Start()
	log.LogInformation("приложение остановлено")
}

func (host *Host) loadPlugins() {
	log := host.loggerFactory.NewLogger("app:plugins")

	host.server.GetRegistry().ClearRoutes()

	path := host.configuration.GetString(cfg.PluginsRootpath)

	if path == nil {
		log.LogCritical("%s was nil", cfg.PluginsRootpath)
		os.Exit(-1)
	}

	err := host.pluginsLoader.LoadPlugins(*path)

	if err != nil {
		log.LogCritical("ошибка загрузки плагинов: %s", err.Error())
		os.Exit(-1)
	}
}

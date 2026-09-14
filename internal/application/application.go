package application

import (
	"os"

	"github.com/sudzekai-web-os/abstractions"
	"github.com/sudzekai-web-os/api/internal/pluginsloader"
)

type application struct {
	server        abstractions.IServer
	loggerFactory abstractions.ILoggerFactory
	executor      abstractions.IExecutor
	pluginsLoader *pluginsloader.PluginsLoader
	configuration abstractions.IConfiguration
}

func (app *application) Run() {
	log := app.loggerFactory.NewLogger("app")

	log.LogInformation("приложение запущено")

	app.loadPlugins()

	app.server.Start()

	log.LogInformation("приложение остановлено")
}

func (app *application) loadPlugins() {
	log := app.loggerFactory.NewLogger("app:plugins")

	app.server.GetRegistry().ClearRoutes()
	err := app.pluginsLoader.LoadPlugins(
		app.configuration.GetString("plugins.rootpath"),
	)

	if err != nil {
		log.LogCritical("ошибка загрузки плагинов: %s", err.Error())
		os.Exit(-1)
	}
}

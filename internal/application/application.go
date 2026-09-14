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

	app.loadModules()

	app.server.Start()

	log.LogInformation("приложение остановлено")
}

func (app *application) loadModules() {
	log := app.loggerFactory.NewLogger("app:modules")

	app.server.GetRegistry().ClearRoutes()
	err := app.pluginsLoader.LoadModules("./modules")

	if err != nil {
		log.LogCritical("ошибка загрузки модулей: %s", err.Error())
		os.Exit(-1)
	}
}

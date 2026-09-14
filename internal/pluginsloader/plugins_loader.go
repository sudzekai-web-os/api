package pluginsloader

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/sudzekai-web-os/abstractions"
)

type PluginsLoader struct {
	plugins        map[string]string
	loggingFactory abstractions.ILoggerFactory
	executor       abstractions.IExecutor
	registry       abstractions.IHandlersRegistry
	logger         abstractions.ILogger
}

func NewPluginsLoader(
	loggerFactory abstractions.ILoggerFactory,
	executor abstractions.IExecutor,
	registry abstractions.IHandlersRegistry) *PluginsLoader {
	return &PluginsLoader{
		plugins:        make(map[string]string),
		loggingFactory: loggerFactory,
		executor:       executor,
		registry:       registry,
		logger:         loggerFactory.NewLogger("plugins-loader"),
	}
}

var plugins map[string]string

func (loader *PluginsLoader) Load(path string) error {
	p, err := plugin.Open(path)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %s", err.Error())
	}

	symbol, err := findSymbol(p)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина:%s", err.Error())
	}

	mod, err := symbolAsModule(symbol)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %s", err.Error())
	}

	if err := mod.Initialize(
		loader.registry,
		loader.loggingFactory,
		loader.executor); err != nil {
		return fmt.Errorf(
			"%s, %s",
			mod.Name(),
			err.Error(),
		)
	}

	plugins[mod.Name()] = mod.Version()

	loader.logger.LogDebug("плагин %s загружен", mod.Name())

	return nil
}

func (loader *PluginsLoader) LoadPlugins(rootPath string) error {
	plugins = make(map[string]string)

	loader.logger.LogInformation("запущена загрузка плагинов...")

	files, err := os.ReadDir(rootPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			loader.logger.LogInformation("загрузка плагина: %s...", file.Name())
			err = loader.Load(filepath.Join("./plugins", file.Name()))
			if err != nil {
				loader.logger.LogError("плагин %s пропущен из-за: %s", file.Name(), err.Error())
			}
		}
	}

	loader.logger.LogInformation("загружено плагинов %d", len(plugins))

	return nil
}

func (loader *PluginsLoader) GetLoadedPlugins() (result []string) {
	result = make([]string, 0)

	for name, ver := range plugins {
		result = append(result, fmt.Sprintf("%s v%s", name, ver))
	}

	return
}

func findSymbol(p *plugin.Plugin) (plugin.Symbol, error) {
	symbol, err := p.Lookup("Module")

	if err != nil {
		return nil, fmt.Errorf("плагин не экспортирует символ Module: %s", err.Error())
	}

	return symbol, nil
}

func symbolAsModule(symbol plugin.Symbol) (abstractions.IModule, error) {
	pluginsPointer, ok := symbol.(*abstractions.IModule)

	if !ok {
		module, ok := symbol.(abstractions.IModule)

		if !ok {
			return nil, fmt.Errorf("плагин не экспортирует корректный тип модуля: Module должен реализовывать интерфейс abstractions.IModule")
		}

		return module, nil
	}

	return *pluginsPointer, nil
}

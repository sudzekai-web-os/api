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
	modules        map[string]string
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
		modules:        make(map[string]string),
		loggingFactory: loggerFactory,
		executor:       executor,
		registry:       registry,
		logger:         loggerFactory.NewLogger("modules-loader"),
	}
}

var modules map[string]string

func (loader *PluginsLoader) Load(path string) error {
	p, err := plugin.Open(path)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %w", err)
	}

	symbol, err := findSymbol(p)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %w", err)
	}

	mod, err := symbolAsModule(symbol)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %w", err)
	}

	if err := mod.Initialize(
		loader.registry,
		loader.loggingFactory,
		loader.executor); err != nil {
		return fmt.Errorf(
			"%s, %w",
			mod.Name(),
			err,
		)
	}

	modules[mod.Name()] = mod.Version()

	loader.logger.LogDebug("плагин %s загружен", mod.Name())

	return nil
}

func (loader *PluginsLoader) LoadModules(rootPath string) error {
	modules = make(map[string]string)

	loader.logger.LogInformation("запущена загрузка плагинов...")

	files, err := os.ReadDir(rootPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			loader.logger.LogInformation("загрузка плагина: %s...", file.Name())
			err = loader.Load(filepath.Join("./modules", file.Name()))
			if err != nil {
				loader.logger.LogError("плагин %s пропущен из-за: %w", file.Name(), err)
			}
		}
	}

	loader.logger.LogInformation("загружено плагинов %d", len(modules))

	return nil
}

func (loader *PluginsLoader) GetLoadedPlugins() (result []string) {
	result = make([]string, 0)

	for name, ver := range modules {
		result = append(result, fmt.Sprintf("%s v%s", name, ver))
	}

	return
}

func findSymbol(p *plugin.Plugin) (plugin.Symbol, error) {
	symbol, err := p.Lookup("Module")

	if err != nil {
		return nil, fmt.Errorf("плагин не экспортирует символ Module: %w", err)
	}

	return symbol, nil
}

func symbolAsModule(symbol plugin.Symbol) (abstractions.IModule, error) {
	modulesPointer, ok := symbol.(*abstractions.IModule)

	if !ok {
		module, ok := symbol.(abstractions.IModule)

		if !ok {
			return nil, fmt.Errorf("плагин не экспортирует корректный тип модуля: Module должен реализовывать интерфейс abstractions.IModule")
		}

		return module, nil
	}

	return *modulesPointer, nil
}

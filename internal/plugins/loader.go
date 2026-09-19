package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"

	"github.com/sudzekai-web-os/core"
)

type PluginsLoader struct {
	plugins       map[string]string
	loggerFactory core.ILoggerFactory
	executor      core.IExecutor
	registry      core.IHandlersRegistry
	logger        core.ILogger
	configuration core.IConfiguration
}

func NewLoader(
	loggerFactory core.ILoggerFactory,
	executor core.IExecutor,
	registry core.IHandlersRegistry,
	configuration core.IConfiguration,
) *PluginsLoader {
	return &PluginsLoader{
		plugins:       make(map[string]string),
		loggerFactory: loggerFactory,
		executor:      executor,
		registry:      registry,
		logger:        loggerFactory.NewLogger("plugins-loader"),
		configuration: configuration,
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
		return fmt.Errorf("ошибка загрузки плагина: %s", err.Error())
	}

	mod, err := symbolAsModule(symbol)

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %s", err.Error())
	}

	loader.TryAddConfiguration(mod)
	loader.TryAddLoggerFactory(mod)
	loader.TryAddHandlersRegistry(mod)
	loader.TryAddExecutor(mod)

	err = mod.Start()

	if err != nil {
		return fmt.Errorf("ошибка загрузки плагина: %s", err.Error())
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

func symbolAsModule(symbol plugin.Symbol) (IPlugin, error) {
	value := reflect.ValueOf(symbol)

	if !value.IsValid() {
		return nil, fmt.Errorf(
			"плагин не экспортирует корректный тип модуля: Module имеет недопустимое значение",
		)
	}

	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, fmt.Errorf(
				"плагин не экспортирует корректный тип модуля: Module равен nil",
			)
		}

		value = value.Elem()
	}

	module, ok := value.Interface().(IPlugin)

	if !ok {
		return nil, fmt.Errorf(
			"плагин не экспортирует корректный тип модуля: Module должен реализовывать интерфейс IPlugin",
		)
	}

	return module, nil
}

func (loader *PluginsLoader) TryAddLoggerFactory(module IPlugin) {
	consumer, ok := module.(ILoggerFactoryConsumer)

	if !ok {
		return
	}

	consumer.AddLoggerFactory(loader.loggerFactory)
}

func (loader *PluginsLoader) TryAddExecutor(module IPlugin) {
	consumer, ok := module.(IExecutorConsumer)

	if !ok {
		return
	}

	consumer.AddExecutor(loader.executor)
}

func (loader *PluginsLoader) TryAddHandlersRegistry(module IPlugin) {
	consumer, ok := module.(IHandlersRegistryConsumer)

	if !ok {
		return
	}

	consumer.AddHandlersRegistry(loader.registry)
}

func (loader *PluginsLoader) TryAddConfiguration(module IPlugin) {
	consumer, ok := module.(IConfigurationConsumer)

	if !ok {
		return
	}

	consumer.AddConfiguration(loader.configuration)
}

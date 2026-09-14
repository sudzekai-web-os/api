package modulesloader

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/sudzekai-web-os/abstractions"
)

type ModulesLoader struct {
	modules        map[string]string
	loggingFactory abstractions.ILoggerFactory
	executor       abstractions.IExecutor
	registry       abstractions.IHandlersRegistry
	logger         abstractions.ILogger
}

func NewModulesLoader(
	loggerFactory abstractions.ILoggerFactory,
	executor abstractions.IExecutor,
	registry abstractions.IHandlersRegistry) *ModulesLoader {
	return &ModulesLoader{
		modules:        make(map[string]string),
		loggingFactory: loggerFactory,
		executor:       executor,
		registry:       registry,
		logger:         loggerFactory.NewLogger("modules-loader"),
	}
}

var modules map[string]string

func (loader *ModulesLoader) Load(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("загрузка модуля %q: %w", path, err)
	}

	symbol, err := p.Lookup("Module")
	if err != nil {
		return fmt.Errorf(
			"модуль %q не экспортирует Module: %w",
			path,
			err,
		)
	}

	modulesPointer, ok := symbol.(*abstractions.IModule)
	if !ok {
		return fmt.Errorf(
			"модуль %q имеет неправильный тип Module",
			path,
		)
	}

	mod := *modulesPointer

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

	loader.logger.LogDebug("модуль %s загружен", mod.Name())

	return nil
}

func (loader *ModulesLoader) LoadModules() error {
	modules = make(map[string]string)

	loader.logger.LogInformation("запущена загрузка модулей...")

	files, err := os.ReadDir("./modules")

	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			loader.logger.LogInformation("загрузка модуля: %s...", file.Name())
			err = loader.Load(filepath.Join("./modules", file.Name()))
			if err != nil {
				loader.logger.LogError("ошибка загрузки модуля: %s. модуль пропущен...", err.Error())
			}
		}
	}

	loader.logger.LogInformation("загружено модулей %d", len(modules))

	return nil
}

func (loader *ModulesLoader) GetLoadedModules() (result []string) {
	result = make([]string, 0)

	for name, ver := range modules {
		result = append(result, fmt.Sprintf("%s v%s", name, ver))
	}

	return
}

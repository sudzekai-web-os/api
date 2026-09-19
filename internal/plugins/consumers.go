package plugins

import "github.com/sudzekai-web-os/core"

type IExecutorConsumer interface {
	AddExecutor(
		executor core.IExecutor,
	)
}

type IHandlersRegistryConsumer interface {
	AddHandlersRegistry(
		registry core.IHandlersRegistry,
	)
}

type ILoggerFactoryConsumer interface {
	AddLoggerFactory(
		loggerFactory core.ILoggerFactory,
	)
}

type IConfigurationConsumer interface {
	AddConfiguration(
		configuration core.IConfiguration,
	)
}

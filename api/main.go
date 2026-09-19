package main

import (
	"strings"

	cfg "github.com/sudzekai-web-os/api/internal/configuration"
	"github.com/sudzekai-web-os/api/internal/objects"
	"github.com/sudzekai-web-os/api/internal/weboscore/configurationbuilder"
	"github.com/sudzekai-web-os/api/internal/weboscore/hostbuilder"
)

var flags = []objects.Flag{
	objects.NewFlag(strings.ReplaceAll(cfg.DatabaseHost, ":", "-"), cfg.DatabaseHost, "localhost", "адрес базы данных"),
	objects.NewFlag(strings.ReplaceAll(cfg.DatabasePassword, ":", "-"), cfg.DatabasePassword, "pwd", "пароль базы данных"),
	objects.NewFlag(strings.ReplaceAll(cfg.DatabasePort, ":", "-"), cfg.DatabasePort, 3306, "порт базы данных"),
	objects.NewFlag(strings.ReplaceAll(cfg.DatabaseSchema, ":", "-"), cfg.DatabaseSchema, "sch", "схема базы данных"),
	objects.NewFlag(strings.ReplaceAll(cfg.DatabaseUser, ":", "-"), cfg.DatabaseUser, "usr", "пользователь базы данных"),

	objects.NewFlag(strings.ReplaceAll(cfg.LoggingLevel, ":", "-"), cfg.LoggingLevel, "information", "уровень логирования"),

	objects.NewFlag(strings.ReplaceAll(cfg.LoggingConsoleEnabled, ":", "-"), cfg.LoggingConsoleEnabled, true, "включить логирование в консоль"),

	objects.NewFlag(strings.ReplaceAll(cfg.LoggingFileEnabled, ":", "-"), cfg.LoggingFileEnabled, true, "включить логирование в файл"),
	objects.NewFlag(strings.ReplaceAll(cfg.LoggingFileName, ":", "-"), cfg.LoggingFileName, "information", "имя файла логирования"),

	objects.NewFlag(strings.ReplaceAll(cfg.LoggingGrpcEnabled, ":", "-"), cfg.LoggingGrpcEnabled, true, "включить логирование grpc"),
	objects.NewFlag(strings.ReplaceAll(cfg.LoggingGrpcEndpoint, ":", "-"), cfg.LoggingGrpcEndpoint, "information", "адрес gRPC"),

	objects.NewFlag(strings.ReplaceAll(cfg.WebhostHost, ":", "-"), cfg.WebhostHost, "0.0.0.0", "адрес HTTP-сервера"),
	objects.NewFlag(strings.ReplaceAll(cfg.WebhostPort, ":", "-"), cfg.WebhostPort, 8080, "порт HTTP-сервера"),

	objects.NewFlag(strings.ReplaceAll(cfg.PluginsRootpath, ":", "-"), cfg.PluginsRootpath, "./plugins", "путь к папке с плагинами (.so)")}

func main() {
	cfgBuilder := configurationbuilder.New()

	cfgBuilder.AddConfigurationFile("configuration.json")
	cfgBuilder.LoadConfigurationFile()

	configureFlags(cfgBuilder)

	cfgBuilder.Parse()

	builder := hostbuilder.New()

	conf := cfgBuilder.GetConfiguration()

	builder.SetConfiguration(conf)

	host := builder.Build()

	host.Run()
}

func configureFlags(cfgBuilder *configurationbuilder.ConfigurationBuilder) {
	for _, flg := range flags {
		cfgBuilder.AddFlag(flg)
	}
}

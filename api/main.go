package main

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/sudzekai/web-os-api/internal/application"
	"github.com/sudzekai/web-os-api/internal/config"
	"github.com/sudzekai/web-os-api/packages/logging"
)

func main() {
	configureEnvironment()

	builder := application.NewApplicationBuilder()

	coonfigureBuilder(builder)

	app, err := builder.Build()

	exitIfErr(err)

	app.Run()
}

func configureEnvironment() {
	configureFlags()
	configureConfig()
}

func configureConfig() {
	pflag.Parse()

	viper.SetConfigName("configuration")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	exitIfErr(err)

	bind := func(key string) {
		if err := viper.BindPFlag(key, pflag.Lookup(strings.ReplaceAll(key, ".", "-"))); err != nil {
			panic(err)
		}
	}

	bind("database.host")
	bind("database.password")
	bind("database.port")
	bind("database.schema")
	bind("database.user")

	bind("logging.default.loglevel")
	bind("logging.system.loglevel")

	bind("webhost.host")
	bind("webhost.port")
}

func configureFlags() {
	pflag.String("database-host", "localhost", "адрес базы данных")
	pflag.String("database-password", "pwd", "пароль базы данных")
	pflag.Int("database-port", 3306, "порт базы данных")
	pflag.String("database-schema", "sch", "схема базы данных")
	pflag.String("database-user", "usr", "пользователь базы данных")

	pflag.String("logging-default-loglevel", "Information", "уровень логирования")
	pflag.String("logging-system-loglevel", "Information", "системный уровень логирования")

	pflag.String("webhost-host", "0.0.0.0", "адрес HTTP-сервера")
	pflag.Int("webhost-port", 8080, "порт HTTP-сервера")
}

func coonfigureBuilder(builder *application.ApplicationBuilder) {
	configuration := config.NewConfiguration(viper.GetViper().AllSettings())

	configuration.SetGetBool(viper.GetBool)
	configuration.SetGetString(viper.GetString)
	configuration.SetGetInt(viper.GetInt)

	err := builder.GetLoggerFactory().SetMinLevelStr(viper.GetString("logging.system.loglevel"))
	exitIfErr(err)

	builder.SetOutput(os.Stdout)
	builder.SetConfiguration(configuration)
}

func exitIfErr(err error) {
	if err != nil {
		log := logging.NewLoggerFactory(os.Stdout).NewLogger("main")
		log.LogCritical("%w", err)
	}
}

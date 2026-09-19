package configurationbuilder

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/sudzekai-web-os/api/internal/objects"
	"github.com/sudzekai-web-os/api/internal/weboscore/configuration"
	"github.com/sudzekai-web-os/core"
)

type ConfigurationBuilder struct {
	fileName    string
	fileOptions map[string]any

	flags []objects.Flag
}

func New() *ConfigurationBuilder {

	builder := &ConfigurationBuilder{
		fileOptions: make(map[string]any),
		flags:       make([]objects.Flag, 0),
	}

	flag.Visit(func(f *flag.Flag) {
		for _, flag := range builder.flags {
			if flag.GetName() == f.Name {
				flag.SetIsVisited(true)
			}
		}
	})

	return builder
}

func (cb *ConfigurationBuilder) AddConfigurationFile(fileName string) *ConfigurationBuilder {
	cb.fileName = fileName
	return cb
}

func (cb *ConfigurationBuilder) AddFlag(newFlag objects.Flag) *ConfigurationBuilder {
	for _, flg := range cb.flags {
		if flg.GetName() == newFlag.GetName() || flg.GetBindProperty() == newFlag.GetBindProperty() {
			panic(fmt.Sprintf(
				"can't add flag: flag with name '%s' or bind property '%s' already exists", newFlag.GetName(), newFlag.GetBindProperty()),
			)
		}
	}

	switch reflect.TypeOf(newFlag.GetValue()) {
	case reflect.TypeFor[int]():
		ptr := flag.Int(newFlag.GetName(), newFlag.GetValue().(int), newFlag.GetUsage())

		newFlag.SetValueGetter(func() any {
			return *ptr
		})

	case reflect.TypeFor[string]():
		ptr := flag.String(newFlag.GetName(), newFlag.GetValue().(string), newFlag.GetUsage())
		newFlag.SetValueGetter(func() any {
			return *ptr
		})

	case reflect.TypeFor[bool]():
		ptr := flag.Bool(newFlag.GetName(), newFlag.GetValue().(bool), newFlag.GetUsage())
		newFlag.SetValueGetter(func() any {
			return *ptr
		})

	default:
		panic(fmt.Sprintf("unknown flag type: %v", reflect.TypeOf(newFlag.GetValue())))
	}

	cb.flags = append(cb.flags, newFlag)

	return cb
}

func (cb *ConfigurationBuilder) Parse() *ConfigurationBuilder {
	flag.Parse()

	return cb
}

func (cb *ConfigurationBuilder) EnableArgs() *ConfigurationBuilder {

	return cb
}

func (cb *ConfigurationBuilder) LoadConfigurationFile() *ConfigurationBuilder {
	if cb.fileName == "" {
		panic("can't load configuration file: name wasn't provided")
	}

	content, err := os.ReadFile(cb.fileName)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &cb.fileOptions)

	if err != nil {
		panic(err)
	}

	return cb
}

func (cb *ConfigurationBuilder) GetConfiguration() core.IConfiguration {
	options := make(map[string]func() any)

	if cb.fileOptions != nil {
		options = makeOptions(cb.fileOptions)
	}

	for _, flg := range cb.flags {
		bindProperty := flg.GetBindProperty()
		_, exists := options[bindProperty]

		if !exists || flg.IsVisited() {
			options[bindProperty] = flg.GetValue
		}
	}

	return configuration.New(options)
}

func makeOptions(config map[string]any) map[string]func() any {
	result := make(map[string]func() any)

	var walk func(map[string]any, string)

	walk = func(config map[string]any, prefix string) {
		for key, value := range config {
			path := key

			if prefix != "" {
				path = prefix + ":" + key
			}

			nested, ok := value.(map[string]any)
			if ok {
				walk(nested, path)
				continue
			}

			if number, ok := value.(float64); ok {
				if number == float64(int(number)) {
					value = int(number)
				}
			}

			result[path] = func() any {
				return value
			}
		}
	}

	walk(config, "")

	return result
}

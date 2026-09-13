package config

import "github.com/sudzekai-web-os/abstractions"

type Configuration struct {
	settings map[string]any

	getString func(key string) string
	getBool   func(key string) bool
	getInt    func(key string) int
}

func NewConfiguration(settings map[string]any) abstractions.IConfiguration {
	return &Configuration{
		settings: settings,
	}
}

func (cfg *Configuration) GetString(key string) string {
	return cfg.getString(key)
}

func (cfg *Configuration) GetBool(key string) bool {
	return cfg.getBool(key)
}

func (cfg *Configuration) GetInt(key string) int {
	return cfg.getInt(key)
}

func (cfg *Configuration) GetSettings() map[string]any {
	return cfg.settings
}

func (cfg *Configuration) SetGetString(fun func(key string) string) {
	cfg.getString = fun
}

func (cfg *Configuration) SetGetBool(fun func(key string) bool) {
	cfg.getBool = fun
}

func (cfg *Configuration) SetGetInt(fun func(key string) int) {
	cfg.getInt = fun
}

package configuration

import "github.com/sudzekai-web-os/core"

type Configuration struct {
	options map[string]func() any
}

func New(options map[string]func() any) core.IConfiguration {
	return Configuration{
		options: options,
	}
}

func (cfg Configuration) GetString(key string) *string {
	val := cfg.GetValue(key)
	if val == nil {
		return nil
	}

	parsed, ok := val.(string)

	if !ok {
		return nil
	}

	return &parsed
}

func (cfg Configuration) GetBool(key string) *bool {
	val := cfg.GetValue(key)
	if val == nil {
		return nil
	}

	parsed, ok := val.(bool)

	if !ok {
		return nil
	}

	return &parsed
}

func (cfg Configuration) GetInt(key string) *int {
	val := cfg.GetValue(key)
	if val == nil {
		return nil
	}

	parsed, ok := val.(int)

	if !ok {
		return nil
	}

	return &parsed
}

func (cfg Configuration) GetValue(key string) any {
	getter, ok := cfg.options[key]
	if !ok {
		return nil
	}

	return getter()
}

func (cfg Configuration) GetOptions() map[string]func() any {
	return cfg.options
}

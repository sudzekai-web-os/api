module github.com/sudzekai/web-os-api

go 1.26.8

require (
	github.com/sudzekai/web-os-api/packages/abstractions v0.0.0 // direct
	github.com/sudzekai/web-os-api/packages/executor v0.0.0 // direct
	github.com/sudzekai/web-os-api/packages/logging v0.0.0 // direct
	github.com/sudzekai/web-os-api/packages/server v0.0.0 // direct
	github.com/sudzekai/web-os-api/packages/types v0.0.0 // direct
	github.com/sudzekai/web-os-api/packages/modulesloader v0.0.0 // direct
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/sudzekai/web-os-api/packages/logging => ../packages/logging
replace github.com/sudzekai/web-os-api/packages/server => ../packages/server
replace github.com/sudzekai/web-os-api/packages/executor => ../packages/executor
replace github.com/sudzekai/web-os-api/packages/abstractions => ../packages/abstractions
replace github.com/sudzekai/web-os-api/packages/types => ../packages/types
replace github.com/sudzekai/web-os-api/packages/modulesloader => ../packages/modulesloader

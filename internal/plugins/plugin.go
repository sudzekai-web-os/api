package plugins

type IPlugin interface {
	Name() string
	Version() string
	Description() string
	Start() error
}

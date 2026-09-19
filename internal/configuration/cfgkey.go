package cfg

const (
	LoggingLevel = "logging:level"

	LoggingConsoleEnabled = "logging:console:enabled"

	LoggingGrpcEnabled = "logging:grpc:enabled"
	LoggingGrpcUrl     = "logging:grpc:url"

	LoggingFileEnabled = "logging:file:enabled"
	LoggingFileName    = "logging:file:name"

	WebhostHost = "webhost:host"
	WebhostPort = "webhost:port"

	DatabaseHost     = "database:host"
	DatabasePassword = "database:password"
	DatabasePort     = "database:port"
	DatabaseSchema   = "database:schema"
	DatabaseUser     = "database:user"

	PluginsRootpath = "plugins:rootpath"
)

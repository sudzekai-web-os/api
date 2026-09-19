package cfg

const (
	LoggingConsoleEnabled = "logging:console:enabled"
	LoggingLevel          = "logging:console:level"

	LoggingGrpcEnabled  = "logging:grpc:enabled"
	LoggingGrpcEndpoint = "logging:grpc:endpoint"

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

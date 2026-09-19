package loggingbuilder

import (
	"os"

	"github.com/sudzekai-web-os/core"
	"github.com/sudzekai-web-os/logging"
)

type LoggingBuilder struct {
	loggerFactory core.ILoggerFactory
}

func New() *LoggingBuilder {
	return &LoggingBuilder{
		loggerFactory: logging.NewLoggerFactory(),
	}
}

func (lb *LoggingBuilder) AddGrpc(endpoint string) *LoggingBuilder {
	lb.loggerFactory.AddWriter(logging.NewGrpcWriter(endpoint))

	return lb
}

func (lb *LoggingBuilder) AddFile(fileName string) *LoggingBuilder {
	lb.loggerFactory.AddWriter(logging.NewFileWriter(fileName))

	return lb
}

func (lb *LoggingBuilder) AddStdout() *LoggingBuilder {
	lb.loggerFactory.AddWriter(logging.NewConsoleWriter(os.Stdout))

	return lb
}

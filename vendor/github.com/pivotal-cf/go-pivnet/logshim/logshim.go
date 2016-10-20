package logshim

import (
	"fmt"
	"log"

	"github.com/pivotal-cf/go-pivnet/logger"
)

type LogShim interface {
	Debug(action string, data ...logger.Data)
	Info(action string, data ...logger.Data)
}

type logShim struct {
	infoLogger  *log.Logger
	debugLogger *log.Logger
	verbose     bool
}

func NewLogShim(
	infoLogger *log.Logger,
	debugLogger *log.Logger,
	verbose bool,
) LogShim {
	return &logShim{
		infoLogger:  infoLogger,
		debugLogger: debugLogger,
		verbose:     verbose,
	}
}

func (l logShim) Debug(action string, data ...logger.Data) {
	if l.verbose {
		l.debugLogger.Println(fmt.Sprintf("%s - %+v", action, data))
	}
}

func (l logShim) Info(action string, data ...logger.Data) {
	l.infoLogger.Println(fmt.Sprintf("%s - %+v", action, data))
}

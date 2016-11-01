package parsehub_go

import (
	"fmt"
	"log"
	"os"
)

type logLevelType int

const (
	LogLevelDebug logLevelType = iota + 1
	LogLevelWarning
	LogLevelFatal
)

var logLevel = LogLevelDebug
var internalLogger *log.Logger


// Set Logger for package
func SetLogger(level logLevelType, logger *log.Logger) {
	internalLogger = logger
	logLevel = level
}

func debugf(s string, args ...interface{}) {
	if internalLogger != nil && logLevel <= LogLevelDebug {
		internalLogger.Output(2, fmt.Sprintf("Debug: " + s, args...))
	}

	return
}

func warningf(s string, args ...interface{}) {
	if internalLogger != nil && logLevel <= LogLevelWarning {
		internalLogger.Output(2, fmt.Sprintf("Warning: " + s, args...))
	}

	return
}

func fatalf(s string, args ...interface{}) {
	if internalLogger != nil && logLevel <= LogLevelFatal {
		internalLogger.Output(2, fmt.Sprintf("Fatal: " + s, args...))
	}
	os.Exit(1)
}

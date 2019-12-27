package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
)

func LogInfo(formatString string, a...  interface{}) {
	processLogMessage(info, formatString, a...)
}

func LogWarning(formatString string, a... interface{}) {
	processLogMessage(warning, formatString, a...)
}

func LogError(formatString string, a... interface{}) {
	processLogMessage(error, formatString, a...)
}

func LogFatal(formatString string, a... interface{}) {
	processLogMessage(fatal, formatString, a...)
}

func VarDumpInfo(message string, object interface{}) {
	processLogMessage(info, message + ": %#v", object)
}

func VarDumpWarning(message string, object interface{}) {
	processLogMessage(warning, message + ": %#v", object)
}

func VarDumpError(message string, object interface{}) {
	processLogMessage(error, message + ": %#v", object)
}

func VarDumpFatal(message string, object interface{}) {
	processLogMessage(fatal, message + ": %#v", object)
}

// Private /////////////////////////////////////////////////////////////////////

type logLevel_t int
const (
	info logLevel_t = iota
	warning
	error
	fatal
)

var logger *log.Logger

/** Package init **/
func init() {
	logger = log.New(os.Stdout, "", 0)
}

func processLogMessage(logLevel logLevel_t, formatString string, a... interface{}) {
	message := fmt.Sprintf(formatString, a...)

	currentTime := time.Now()
	prefix := currentTime.Format("2006-01-02 15:04:05.000000") + " | "
	suffix := ""

	switch logLevel {
	case info:
		prefix += "LogLevel=info | "
	case warning:
		prefix += "LogLevel=warning | "
	case error:
		prefix += "LogLevel=error | "
		suffix += "\nStacktrace: " + string(debug.Stack())
	case fatal:
		prefix += "LogLevel=fatal | "
		suffix += "\nStacktrace: " + string(debug.Stack())
	}

	logger.Println(prefix + message + suffix)
	if logLevel == fatal {
		panic("Shutting down due to log fatal: " + message)
	}
}

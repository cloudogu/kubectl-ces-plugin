package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloudogu/cesapp-lib/core"
)

var (
	loggerInstance *NamedLogger
	formatter      = createFormatter()
)

func init() {
	initLoggerWithThrowAway()
}

func initLoggerWithThrowAway() {
	// the loggers-to-be-used are instantiated later on when the CLI log level flags have been parsed
	log := logrus.New()
	namedLog := &NamedLogger{logger: log, name: "throwAway"}
	loggerInstance = namedLog
	core.GetLogger = func() core.Logger {
		return namedLog
	}
}

func GetInstance() *NamedLogger {
	return loggerInstance
}

func createFormatter() *Formatter {
	formatter := &Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: time.RFC3339,
		NoColors:        true,
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			file := strings.Split(frame.File, "/")
			function := strings.Split(frame.Function, ".")
			return file[len(file)-1] + ":" + function[len(function)-1] + ":" + strconv.Itoa(frame.Line) + " "
		},
	}

	return formatter
}

// NamedLogger provides log facilities that can be replaced during runtime.
type NamedLogger struct {
	logger logger
	name   string
}

func (ll *NamedLogger) log(level int, args ...interface{}) {
	ll.logger.Info(level, fmt.Sprintf("[%s] %s", ll.name, fmt.Sprint(args...)))
}

func (ll *NamedLogger) logf(level int, format string, args ...interface{}) {
	ll.logger.Info(level, fmt.Sprintf("[%s] %s", ll.name, fmt.Sprintf(format, args...)))
}

// Debug logs a message with debug level.
func (ll *NamedLogger) Debug(args ...interface{}) {
	ll.log(debugLevel, args...)
}

// Info logs a message with info level.
func (ll *NamedLogger) Info(args ...interface{}) {
	ll.log(infoLevel, args...)
}

// Warning logs a message with warn level.
func (ll *NamedLogger) Warning(args ...interface{}) {
	ll.log(warningLevel, args...)
}

// Error logs a message with error level.
func (ll *NamedLogger) Error(args ...interface{}) {
	ll.log(errorLevel, args...)
}

// Debugf logs a printf-style message with debug level.
func (ll *NamedLogger) Debugf(format string, args ...interface{}) {
	ll.logf(debugLevel, format, args...)
}

// Infof logs a printf-style message with info level.
func (ll *NamedLogger) Infof(format string, args ...interface{}) {
	ll.logf(infoLevel, format, args...)
}

// Warningf logs a printf-style message with warn level.
func (ll *NamedLogger) Warningf(format string, args ...interface{}) {
	ll.logf(warningLevel, format, args...)
}

// Errorf logs a printf-style message with error level.
func (ll *NamedLogger) Errorf(format string, args ...interface{}) {
	ll.logf(errorLevel, format, args...)
}

// ConfigureLogger sets up both the logging framework native to this application as well to used libraries.
func ConfigureLogger() error {
	level, err := getLogLevelFromEnv()
	if err != nil {
		return fmt.Errorf("could not configure log level")
	}

	logrusLog := logrus.New()
	logrusLog.SetFormatter(formatter)
	logrusLog.SetLevel(level)

	loggerInstance = &NamedLogger{name: "kubectl-ces", logger: logrusLog}

	cesappLibLogger := &NamedLogger{name: "cesapp-lib", logger: logrusLog}
	// overwrite cesapp-lib logger with our own custom logger implementation to gain logs from cesapp-lib
	core.GetLogger = func() core.Logger {
		return cesappLibLogger
	}

	return nil
}

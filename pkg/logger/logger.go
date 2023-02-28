package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

const LogLevelKey = "log-level"

const (
	LogLevelError LogLevel = "error"
	LogLevelWarn  LogLevel = "warn"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

// logrus log level
const (
	errorLevel int = iota
	warningLevel
	infoLevel
	debugLevel
)

var (
	loggerInstance *CesappLogger
	formatter      = createFormatter()
)

type CesappLogger struct {
	logger *logrus.Logger
}

func init() {
	loggerInstance = NewCesappLogger()
	core.GetLogger = func() core.Logger {
		return GetInstance()
	}
}

func GetInstance() logrus.FieldLogger {
	return loggerInstance
}

func NewCesappLogger() *CesappLogger {
	throwAwayLogger := logrus.New()
	return &CesappLogger{logger: throwAwayLogger}
}

func createFormatter() Formatter {
	formatter := Formatter{
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

func getLogLevelFromEnv() logrus.Level {
	logLevelRef := viper.GetViper().Get(util.CliTransportLogLevel)
	logLevel := *(logLevelRef).(*LogLevel)

	switch logLevel {
	case LogLevelError:
		return logrus.ErrorLevel
	case LogLevelWarn:
		return logrus.WarnLevel
	case LogLevelInfo:
		return logrus.InfoLevel
	case LogLevelDebug:
		return logrus.DebugLevel
	default:
		panic("unsupported log level should be caught during the CLI flag parsing")
	}
}

type LogLevel string

// String returns the stringn value of a log level.
func (ll *LogLevel) String() string {
	return string(*ll)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (ll *LogLevel) Set(v string) error {
	switch v {
	case string(LogLevelError), string(LogLevelWarn), string(LogLevelInfo), string(LogLevelDebug):
		*ll = LogLevel(v)
		return nil
	default:
		return fmt.Errorf(`log-level must be one of "%s", "%s", "%s", or "%s"`, LogLevelError, LogLevelWarn, LogLevelInfo, LogLevelDebug)
	}
}

// Type is only used in help text
func (ll *LogLevel) Type() string {
	return LogLevelKey
}

type libraryLogger struct {
	logger *logrus.Logger
	name   string
}

func (ll *libraryLogger) log(level int, args ...interface{}) {
	ll.logger.Info(level, fmt.Sprintf("[%s] %s", ll.name, fmt.Sprint(args...)))
}

func (ll *libraryLogger) logf(level int, format string, args ...interface{}) {
	ll.logger.Info(level, fmt.Sprintf("[%s] %s", ll.name, fmt.Sprintf(format, args...)))
}

func (ll *libraryLogger) Debug(args ...interface{}) {
	ll.log(debugLevel, args...)
}

func (ll *libraryLogger) Info(args ...interface{}) {
	ll.log(infoLevel, args...)
}

func (ll *libraryLogger) Warning(args ...interface{}) {
	ll.log(warningLevel, args...)
}

func (ll *libraryLogger) Error(args ...interface{}) {
	ll.log(errorLevel, args...)
}

func (ll *libraryLogger) Debugf(format string, args ...interface{}) {
	ll.logf(debugLevel, format, args...)
}

func (ll *libraryLogger) Infof(format string, args ...interface{}) {
	ll.logf(infoLevel, format, args...)
}

func (ll *libraryLogger) Warningf(format string, args ...interface{}) {
	ll.logf(warningLevel, format, args...)
}

func (ll *libraryLogger) Errorf(format string, args ...interface{}) {
	ll.logf(errorLevel, format, args...)
}

func ConfigureLogger() {
	level := getLogLevelFromEnv()

	// create logrus logger that can be styled and formatted
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.TextFormatter{})
	logrusLog.SetLevel(level)

	// set custom logger implementation to cesapp-lib logger
	cesappLibLogger := libraryLogger{name: "cesapp-lib", logger: logrusLog}
	core.GetLogger = func() core.Logger {
		return &cesappLibLogger
	}
}

func (c *CesappLogger) WithField(key string, value interface{}) *logrus.Entry {
	return c.logger.WithField(key, value)
}

func (c *CesappLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return c.logger.WithFields(fields)
}

func (c *CesappLogger) WithError(err error) *logrus.Entry {
	return c.logger.WithError(err)
}

func (c *CesappLogger) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}

func (c *CesappLogger) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}

func (c *CesappLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	c.logger.Tracef(format, args...)
}

func (c *CesappLogger) Warnf(format string, args ...interface{}) {
	c.logger.Warnf(format, args...)
}

func (c *CesappLogger) Warningf(format string, args ...interface{}) {
	c.logger.Warningf(format, args...)
}

func (c *CesappLogger) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}

func (c *CesappLogger) Fatalf(format string, args ...interface{}) {
	c.logger.Fatalf(format, args...)
}

func (c *CesappLogger) Panicf(format string, args ...interface{}) {
	c.logger.Panicf(format, args...)
}

func (c *CesappLogger) Debug(args ...interface{}) {
	c.logger.Debug(args...)
}

func (c *CesappLogger) Info(args ...interface{}) {
	c.logger.Info(args...)
}

func (c *CesappLogger) Print(args ...interface{}) {
	fmt.Print(args...)
	c.logger.Trace(args...)
}

func (c *CesappLogger) Warn(args ...interface{}) {
	c.logger.Warn(args...)
}

func (c *CesappLogger) Warning(args ...interface{}) {
	c.logger.Warning(args...)
}

func (c *CesappLogger) Error(args ...interface{}) {
	c.logger.Error(args...)
}

func (c *CesappLogger) Fatal(args ...interface{}) {
	c.logger.Fatal(args...)
}

func (c *CesappLogger) Panic(args ...interface{}) {
	c.logger.Panic(args...)
}

func (c *CesappLogger) Debugln(args ...interface{}) {
	c.logger.Debugln(args...)
}

func (c *CesappLogger) Infoln(args ...interface{}) {
	c.logger.Infoln(args...)
}

func (c *CesappLogger) Println(args ...interface{}) {
	fmt.Println(args...)
	c.logger.Traceln(args...)
}

func (c *CesappLogger) Warnln(args ...interface{}) {
	c.logger.Warnln(args...)
}

func (c *CesappLogger) Warningln(args ...interface{}) {
	c.logger.Warningln(args...)
}

func (c *CesappLogger) Errorln(args ...interface{}) {
	c.logger.Errorln(args...)
}

func (c *CesappLogger) Fatalln(args ...interface{}) {
	c.logger.Fatalln(args...)
}

func (c *CesappLogger) Panicln(args ...interface{}) {
	c.logger.Panicln(args...)
}

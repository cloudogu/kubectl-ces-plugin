package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

// LogLevelKey identifies the log-level type.
const LogLevelKey = "log-level"

// CLI log level strings
const (
	LogLevelError LogLevel = "error"
	LogLevelWarn  LogLevel = "warn"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

// logrus log levels
const (
	errorLevel int = iota
	warningLevel
	infoLevel
	debugLevel
)

func getLogLevelFromEnv() (logrus.Level, error) {
	logLevelRef := viper.GetViper().Get(util.CliTransportLogLevel)
	logLevel, ok := logLevelRef.(*LogLevel)
	if !ok {
		return logrus.WarnLevel, fmt.Errorf("log level is not of expected type")
	}

	switch *logLevel {
	case LogLevelError:
		return logrus.ErrorLevel, nil
	case LogLevelWarn:
		return logrus.WarnLevel, nil
	case LogLevelInfo:
		return logrus.InfoLevel, nil
	case LogLevelDebug:
		return logrus.DebugLevel, nil
	default:
		panic("unsupported log level should be caught during the CLI flag parsing")
	}
}

// LogLevel implements a CLI flag.
type LogLevel string

// String returns the string value of a log level flag.
func (ll *LogLevel) String() string {
	return string(*ll)
}

// Set stores the given flag value.
func (ll *LogLevel) Set(v string) error {
	switch v {
	case string(LogLevelError), string(LogLevelWarn), string(LogLevelInfo), string(LogLevelDebug):
		*ll = LogLevel(v)
		return nil
	default:
		return fmt.Errorf(`log-level must be one of "%s", "%s", "%s", or "%s"`, LogLevelError, LogLevelWarn, LogLevelInfo, LogLevelDebug)
	}
}

// Type returns the type of the CLI flag
func (ll *LogLevel) Type() string {
	return LogLevelKey
}

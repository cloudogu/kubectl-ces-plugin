package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

type Logger struct {
}

func NewLogger() *Logger {
	logLevelRef := viper.GetViper().Get(util.CliTransportLogLevel)
	logLevel := *(logLevelRef).(*LogLevel)
	fmt.Println(logLevel)
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if msg == "" {
		_, _ = fmt.Fprintln(os.Stdout, "")
		return
	}

	c := color.New(color.FgHiCyan)
	_, _ = c.Fprintln(os.Stdout, fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(err error) {
	c := color.New(color.FgHiRed)
	_, _ = c.Fprintln(os.Stderr, fmt.Sprintf("%#v", err))
}

func (l *Logger) Instructions(msg string, args ...interface{}) {
	white := color.New(color.FgHiWhite)
	_, _ = white.Fprintln(os.Stdout, "")
	_, _ = white.Fprintln(os.Stdout, fmt.Sprintf(msg, args...))
}

type LogLevel string

const LogLevelKey = "log-level"

const (
	LogLevelError LogLevel = "error"
	LogLevelWarn  LogLevel = "warn"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

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

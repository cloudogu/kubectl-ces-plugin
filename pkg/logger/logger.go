package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Logger struct {
}

func NewLogger() *Logger {
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

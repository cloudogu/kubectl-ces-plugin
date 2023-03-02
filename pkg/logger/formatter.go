package logger

import (
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	optionalCIPathPattern              = "(_.+)?"   // allow paths like ".../<pkg>_PR-2/..."
	optionalVersionedModulePathPattern = `(v\d+/)?` // allow paths like ".../<pkg>/v2/..."
)

var loggerFilePathPattern = fmt.Sprintf("kubectl-ces-plugin%s/pkg/%slogger/logger.go", optionalCIPathPattern, optionalVersionedModulePathPattern)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// CallerFirst - print caller info first
	CallerFirst bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	// output buffer
	b := &bytes.Buffer{}

	f.writeTimeAndLevel(entry, b)

	err := f.writeFieldsAndMessage(entry, b)
	if err != nil {
		return nil, err
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeTimeAndLevel(entry *logrus.Entry, b *bytes.Buffer) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if !f.NoColors {
		_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}
}

func (f *Formatter) writeFieldsAndMessage(entry *logrus.Entry, b *bytes.Buffer) error {
	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	if f.CallerFirst {
		err := f.writeCaller(b, entry)
		if err != nil {
			return fmt.Errorf("failed to write caller: %w", err)
		}
	}

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if !f.CallerFirst {
		err := f.writeCaller(b, entry)
		if err != nil {
			return fmt.Errorf("failed to write caller: %w", err)
		}
	}
	return nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) error {
	if entry.HasCaller() {
		pathExpression, err := regexp.Compile(loggerFilePathPattern)
		if err != nil {
			return errors.Wrap(err, "failed to compile regular expression for module version path")
		}
		caller := getFrame(*pathExpression)
		if f.CustomCallerFormatter != nil {
			_, _ = fmt.Fprint(b, f.CustomCallerFormatter(caller))
		} else {
			_, _ = fmt.Fprintf(
				b,
				" (%s:%d %s)",
				caller.File,
				caller.Line,
				caller.Function,
			)
		}
	}
	return nil
}

func getFrame(loggerFileNameExpression regexp.Regexp) *runtime.Frame {
	fallBackExitAtFrameIndex := 20 // the original caller index is about 13 frames deep but exit 20 frames at the latest

	programCounters := make([]uintptr, fallBackExitAtFrameIndex)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		logDelegateMethodFound := false
		for more, frameIndex := true, 0; more && frameIndex <= fallBackExitAtFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()

			// searching our log delegate method is the fix point that shows us the relation to the original caller
			if !logDelegateMethodFound {
				if foundMatch := loggerFileNameExpression.MatchString(frameCandidate.File); foundMatch {
					logDelegateMethodFound = true
				}
				continue
			}

			if logDelegateMethodFound {
				// the previous "continue" call found the frame before the original caller. So this frame must contain the original caller
				frame = frameCandidate
				break
			}
		}
	}

	return &frame
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if !foundFieldsMap[field] {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		_, _ = fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		_, _ = fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

const (
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorGreen
	}
}

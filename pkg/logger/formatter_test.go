package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_getFrame(t *testing.T) {
	testDelegateMethodeOneFrameBeforeThisFramesTest := fmt.Sprintf("kubectl-ces-plugin%s/pkg/%slogger/formatter.go", optionalCIPathPattern, optionalVersionedModulePathPattern)
	filePathExpression, err := regexp.Compile(testDelegateMethodeOneFrameBeforeThisFramesTest)
	require.NoError(t, err)

	// when
	frameForThisLine := getFrame(*filePathExpression)

	// then
	require.Regexp(t, regexp.MustCompile("kubectl-ces-plugin(_.+)?/pkg/logger/formatter_test.go"), frameForThisLine.File)
	// asserting a specific line in source code makes this test error prone
	// but it's really unlikely that the line number is zero if the function did right
	require.NotEmpty(t, frameForThisLine.Line)
}

func Test_loggerFilePathPattern(t *testing.T) {
	testExpression, err := regexp.Compile(loggerFilePathPattern)
	require.NoError(t, err)
	t.Run("should match path if used in kubectl-ces-plugin", func(t *testing.T) {
		testPath := "kubectl-ces-plugin_PR-2/pkg/logger/logger.go"

		foundMatch := testExpression.MatchString(testPath)

		require.True(t, foundMatch)
	})
}

func Test_getColorByLevel(t *testing.T) {
	tests := []struct {
		name  string
		level logrus.Level
		want  int
	}{
		{
			name:  "should default to green",
			level: math.MaxInt32,
			want:  colorGreen,
		},
		{
			name:  "should be gray on debug",
			level: logrus.DebugLevel,
			want:  colorGray,
		},
		{
			name:  "should be yellow on warn",
			level: logrus.WarnLevel,
			want:  colorYellow,
		},
		{
			name:  "should be red on error",
			level: logrus.ErrorLevel,
			want:  colorRed,
		},
		{
			name:  "should be red on fatal",
			level: logrus.FatalLevel,
			want:  colorRed,
		},
		{
			name:  "should be red on panic",
			level: logrus.PanicLevel,
			want:  colorRed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColorByLevel(tt.level); got != tt.want {
				t.Errorf("getColorByLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_writeField(t *testing.T) {
	type fields struct {
		HideKeys      bool
		NoFieldsSpace bool
	}
	type args struct {
		entry *logrus.Entry
		field string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected string
	}{
		{
			name: "should hide keys with no trailing space",
			fields: fields{
				HideKeys:      true,
				NoFieldsSpace: true,
			},
			args: args{
				entry: &logrus.Entry{
					Data: logrus.Fields{
						"key1": "value1",
						"key2": "value2",
					},
				},
				field: "key2",
			},
			expected: "[value2]",
		},
		{
			name: "should show keys with trailing space",
			fields: fields{
				HideKeys:      false,
				NoFieldsSpace: false,
			},
			args: args{
				entry: &logrus.Entry{
					Data: logrus.Fields{
						"key1": "value1",
						"key2": "value2",
					},
				},
				field: "key2",
			},
			expected: "[key2:value2] ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			b := &bytes.Buffer{}
			f := &Formatter{
				HideKeys:      tt.fields.HideKeys,
				NoFieldsSpace: tt.fields.NoFieldsSpace,
			}

			// when
			f.writeField(b, tt.args.entry, tt.args.field)

			// then
			assert.Equal(t, tt.expected, b.String())
		})
	}
}

func TestFormatter_writeOrderedFields(t *testing.T) {
	tests := []struct {
		name        string
		fieldsOrder []string
		entry       *logrus.Entry
		expected    string
	}{
		{
			name:        "should write ordered fields",
			fieldsOrder: []string{"a", "b", "c"},
			entry: &logrus.Entry{
				Data: logrus.Fields{
					"b": "banana",
					"a": "ananas",
					"c": "coconut",
				},
			},
			expected: "[a:ananas] [b:banana] [c:coconut] ",
		},
		{
			name:        "should write ordered fields with some fields not found",
			fieldsOrder: []string{"a", "b", "c"},
			entry: &logrus.Entry{
				Data: logrus.Fields{
					"b": "banana",
					"d": "durian",
					"a": "ananas",
					"c": "coconut",
				},
			},
			expected: "[a:ananas] [b:banana] [c:coconut] [d:durian] ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			b := &bytes.Buffer{}
			f := &Formatter{
				FieldsOrder: tt.fieldsOrder,
			}

			// when
			f.writeOrderedFields(b, tt.entry)

			// then
			assert.Equal(t, tt.expected, b.String())
		})
	}
}

func TestFormatter_writeFields(t *testing.T) {
	tests := []struct {
		name     string
		entry    *logrus.Entry
		expected string
	}{
		{
			name: "should write nothing with no fields",
			entry: &logrus.Entry{
				Data: logrus.Fields{},
			},
			expected: "",
		},
		{
			name: "should write fields alphabetically sorted",
			entry: &logrus.Entry{
				Data: logrus.Fields{
					"b": "banana",
					"a": "ananas",
					"c": "coconut",
				},
			},
			expected: "[a:ananas] [b:banana] [c:coconut] ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			b := &bytes.Buffer{}
			f := &Formatter{}

			// when
			f.writeFields(b, tt.entry)

			// then
			assert.Equal(t, tt.expected, b.String())
		})
	}
}

func TestFormatter_writeCaller(t *testing.T) {
	originalLoggerFilePathPattern := loggerFilePathPattern
	tests := []struct {
		name                  string
		loggerFilePathPattern string
		customCallerFormatter func(*runtime.Frame) string
		entry                 *logrus.Entry
		expected              string
		wantErr               func(t *testing.T, err error)
	}{
		{
			name: "should do nothing if no caller exists",
			entry: &logrus.Entry{
				Caller: nil,
			},
			expected: "",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:                  "should fail to compile regex",
			loggerFilePathPattern: "(()",
			entry: &logrus.Entry{
				Logger: &logrus.Logger{
					ReportCaller: true,
				},
				Caller: &runtime.Frame{},
			},
			expected: "",
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "error parsing regexp")
			},
		},
		{
			name: "should write without custom caller formatter",
			entry: &logrus.Entry{
				Logger: &logrus.Logger{
					ReportCaller: true,
				},
				Caller: &runtime.Frame{},
			},
			expected: " (:0 unknown)",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "should write with custom caller formatter",
			customCallerFormatter: func(_ *runtime.Frame) string {
				return "mango"
			},
			entry: &logrus.Entry{
				Logger: &logrus.Logger{
					ReportCaller: true,
				},
				Caller: &runtime.Frame{},
			},
			expected: "mango",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			defer func() { loggerFilePathPattern = originalLoggerFilePathPattern }()
			if tt.loggerFilePathPattern != "" {
				loggerFilePathPattern = tt.loggerFilePathPattern
			}
			b := &bytes.Buffer{}
			f := &Formatter{
				CustomCallerFormatter: tt.customCallerFormatter,
			}

			// when
			err := f.writeCaller(b, tt.entry)

			// then
			tt.wantErr(t, err)
			assert.Equal(t, tt.expected, b.String())
		})
	}
}

func TestFormatter_Format(t *testing.T) {
	type fields struct {
		FieldsOrder           []string
		TimestampFormat       string
		HideKeys              bool
		NoColors              bool
		NoFieldsColors        bool
		NoFieldsSpace         bool
		ShowFullLevel         bool
		NoUppercaseLevel      bool
		TrimMessages          bool
		CallerFirst           bool
		CustomCallerFormatter func(*runtime.Frame) string
	}
	tests := []struct {
		name    string
		fields  fields
		arg     *logrus.Entry
		want    string
		wantErr func(t *testing.T, err error)
	}{
		{
			name: "should succeed",
			fields: fields{
				FieldsOrder:      []string{"a", "b", "c"},
				TimestampFormat:  time.RFC3339,
				HideKeys:         false,
				NoColors:         false,
				NoFieldsColors:   false,
				NoFieldsSpace:    false,
				ShowFullLevel:    false,
				NoUppercaseLevel: false,
				TrimMessages:     false,
				CallerFirst:      false,
				CustomCallerFormatter: func(_ *runtime.Frame) string {
					return " (caller) "
				},
			},
			arg: &logrus.Entry{
				Message: "error in my fruity program",
				Level:   logrus.ErrorLevel,
				Data: logrus.Fields{
					"b": "banana",
					"d": "durian",
					"a": "ananas",
					"c": "coconut",
				},
				Logger: &logrus.Logger{ReportCaller: true},
				Caller: &runtime.Frame{},
			},
			want: "0001-01-01T00:00:00Z\x1b[31m [ERRO] [a:ananas] [b:banana] [c:coconut] [d:durian] \x1b[0merror in my fruity program (caller) \n",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "should succeed with different settings",
			fields: fields{
				HideKeys:         true,
				NoColors:         true,
				NoFieldsColors:   true,
				NoFieldsSpace:    true,
				ShowFullLevel:    true,
				NoUppercaseLevel: true,
				TrimMessages:     true,
				CallerFirst:      true,
				CustomCallerFormatter: func(_ *runtime.Frame) string {
					return " (caller) "
				},
			},
			arg: &logrus.Entry{
				Message: "error in my fruity program",
				Level:   logrus.ErrorLevel,
				Data: logrus.Fields{
					"b": "banana",
					"d": "durian",
					"a": "ananas",
					"c": "coconut",
				},
				Logger: &logrus.Logger{ReportCaller: true},
				Caller: &runtime.Frame{},
			},
			want: "Jan  1 00:00:00.000 [error][ananas][banana][coconut][durian]  (caller) error in my fruity program\n",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			f := &Formatter{
				FieldsOrder:           tt.fields.FieldsOrder,
				TimestampFormat:       tt.fields.TimestampFormat,
				HideKeys:              tt.fields.HideKeys,
				NoColors:              tt.fields.NoColors,
				NoFieldsColors:        tt.fields.NoFieldsColors,
				NoFieldsSpace:         tt.fields.NoFieldsSpace,
				ShowFullLevel:         tt.fields.ShowFullLevel,
				NoUppercaseLevel:      tt.fields.NoUppercaseLevel,
				TrimMessages:          tt.fields.TrimMessages,
				CallerFirst:           tt.fields.CallerFirst,
				CustomCallerFormatter: tt.fields.CustomCallerFormatter,
			}

			// when
			got, err := f.Format(tt.arg)

			// then
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

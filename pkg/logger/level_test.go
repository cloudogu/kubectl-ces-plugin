package logger

import (
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getLogLevelFromEnv(t *testing.T) {
	originalLogLevel := viper.GetViper().Get(util.CliTransportLogLevel)

	logLevelError := LogLevelError
	logLevelWarn := LogLevelWarn
	logLevelInfo := LogLevelInfo
	logLevelDebug := LogLevelDebug
	type test struct {
		name     string
		logLevel interface{}
		want     logrus.Level
		wantErr  func(t *testing.T, err error)
	}
	tests := []test{
		{
			name:     "should fail due to unexpected type",
			logLevel: "unexpected type",
			want:     logrus.WarnLevel,
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "log level is not of expected type")
			},
		},
		{
			name:     "should return error level",
			logLevel: &logLevelError,
			want:     logrus.ErrorLevel,
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "should return error level",
			logLevel: &logLevelWarn,
			want:     logrus.WarnLevel,
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "should return error level",
			logLevel: &logLevelInfo,
			want:     logrus.InfoLevel,
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "should return error level",
			logLevel: &logLevelDebug,
			want:     logrus.DebugLevel,
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			defer viper.GetViper().Set(util.CliTransportLogLevel, originalLogLevel)
			viper.GetViper().Set(util.CliTransportLogLevel, tt.logLevel)

			// when
			got, err := getLogLevelFromEnv()

			// then
			tt.wantErr(t, err)
			assert.Equalf(t, tt.want, got, "getLogLevelFromEnv()")
		})
	}

	t.Run("should panic on unsupported log level", func(t *testing.T) {
		// given
		defer viper.GetViper().Set(util.CliTransportLogLevel, originalLogLevel)
		var logLevel LogLevel = "unsupported"
		viper.GetViper().Set(util.CliTransportLogLevel, &logLevel)

		// when
		assert.PanicsWithValue(t, "unsupported log level should be caught during the CLI flag parsing", func() {
			_, _ = getLogLevelFromEnv()
		})
	})
}

func TestLogLevel_String(t *testing.T) {
	t.Run("should return correct string", func(t *testing.T) {
		// given
		sut := LogLevelWarn

		// when
		actual := sut.String()

		// then
		assert.Equal(t, "warn", actual)
	})
}

func TestLogLevel_Set(t *testing.T) {
	emptyLogLevel := LogLevel("")
	tests := []struct {
		name      string
		sut       *LogLevel
		flagValue string
		want      string
		wantErr   func(t *testing.T, err error)
	}{
		{
			name:      "should set error level",
			sut:       &emptyLogLevel,
			flagValue: "error",
			want:      "error",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "should set warn level",
			sut:       &emptyLogLevel,
			flagValue: "warn",
			want:      "warn",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "should set info level",
			sut:       &emptyLogLevel,
			flagValue: "info",
			want:      "info",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "should set debug level",
			sut:       &emptyLogLevel,
			flagValue: "debug",
			want:      "debug",
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "should fail on unsupported log level",
			sut:       &emptyLogLevel,
			flagValue: "unsupported",
			want:      "",
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "log-level must be one of \"error\", \"warn\", \"info\", or \"debug\"")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			*tt.sut = ""

			// when
			err := tt.sut.Set(tt.flagValue)

			// then
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, tt.sut.String())
		})
	}
}

func TestLogLevel_Type(t *testing.T) {
	t.Run("should return correct key", func(t *testing.T) {
		// given
		sut := LogLevelWarn

		// when
		actual := sut.Type()

		// then
		assert.Equal(t, "log-level", actual)
	})
}

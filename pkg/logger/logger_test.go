package logger

import (
	"fmt"
	"testing"
)

func TestNamedLogger_Debug(t *testing.T) {
	t.Run("should print debug log", func(t *testing.T) {
		// given
		const args = "test-output"
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(debugLevel, fmt.Sprintf("[test-logger] %s", args))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Debug(args)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Info(t *testing.T) {
	t.Run("should print info log", func(t *testing.T) {
		// given
		const args = "test-output"
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(infoLevel, fmt.Sprintf("[test-logger] %s", args))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Info(args)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Warning(t *testing.T) {
	t.Run("should print warning log", func(t *testing.T) {
		// given
		const args = "test-output"
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(warningLevel, fmt.Sprintf("[test-logger] %s", args))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Warning(args)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Error(t *testing.T) {
	t.Run("should print error log", func(t *testing.T) {
		// given
		const args = "test-output"
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(errorLevel, fmt.Sprintf("[test-logger] %s", args))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Error(args)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Debugf(t *testing.T) {
	t.Run("should print formatted debug log", func(t *testing.T) {
		// given
		const formatString = "test-output (%s) (%d)"
		const formatArg1 = "some string"
		const formatArg2 = 42
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(debugLevel, fmt.Sprintf("[test-logger] %s", fmt.Sprintf(formatString, formatArg1, formatArg2)))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Debugf(formatString, formatArg1, formatArg2)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Infof(t *testing.T) {
	t.Run("should print formatted info log", func(t *testing.T) {
		// given
		const formatString = "test-output (%s) (%d)"
		const formatArg1 = "some string"
		const formatArg2 = 42
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(infoLevel, fmt.Sprintf("[test-logger] %s", fmt.Sprintf(formatString, formatArg1, formatArg2)))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Infof(formatString, formatArg1, formatArg2)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Warningf(t *testing.T) {
	t.Run("should print formatted warning log", func(t *testing.T) {
		// given
		const formatString = "test-output (%s) (%d)"
		const formatArg1 = "some string"
		const formatArg2 = 42
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(warningLevel, fmt.Sprintf("[test-logger] %s", fmt.Sprintf(formatString, formatArg1, formatArg2)))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Warningf(formatString, formatArg1, formatArg2)

		// then
		// asserts done by mocks
	})
}

func TestNamedLogger_Errorf(t *testing.T) {
	t.Run("should print formatted error log", func(t *testing.T) {
		// given
		const formatString = "test-output (%s) (%d)"
		const formatArg1 = "some string"
		const formatArg2 = 42
		loggerMock := newMockLogger(t)
		loggerMock.EXPECT().Info(errorLevel, fmt.Sprintf("[test-logger] %s", fmt.Sprintf(formatString, formatArg1, formatArg2)))
		sut := &NamedLogger{
			logger: loggerMock,
			name:   "test-logger",
		}

		// when
		sut.Errorf(formatString, formatArg1, formatArg2)

		// then
		// asserts done by mocks
	})
}

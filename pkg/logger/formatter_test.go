package logger

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getFrame(t *testing.T) {
	testDelegateMethodeOneFrameBeforeThisFramesTest := fmt.Sprintf("cesapp/%slogging/formatter.go", modulePathPattern)
	filePathExpression, err := regexp.Compile(testDelegateMethodeOneFrameBeforeThisFramesTest)
	require.NoError(t, err)

	// when
	frameForThisLine := getFrame(*filePathExpression)

	// then
	require.Contains(t, frameForThisLine.File, "cesapp/logging/formatter_test.go")
	// asserting a specific line in source code makes this test error prone
	// but it's really unlikely that the line number is zero if the function did right
	require.NotEmpty(t, frameForThisLine.Line)
}

func Test_loggerFilePathPattern(t *testing.T) {
	testExpression, err := regexp.Compile(loggerFilePathPattern)
	require.NoError(t, err)
	t.Run("should match path if used in cesapp", func(t *testing.T) {
		testPath := "cesapp/logging/logger.go"

		foundMatch := testExpression.MatchString(testPath)

		require.True(t, foundMatch)
	})

	t.Run("should match path if used as module", func(t *testing.T) {
		testPath := "cesapp/v2/logging/logger.go"

		foundMatch := testExpression.MatchString(testPath)

		require.True(t, foundMatch)
	})
}

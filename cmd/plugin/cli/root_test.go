package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getAllForDoguCmd(t *testing.T) {

	actual := new(bytes.Buffer)
	RootCmd := RootCmd()
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs([]string{"dogu", "config", "get", "redmine"})
	err := RootCmd.Execute()

	assert.NoError(t, err, "command should be successful")
	//TODO: test specific output
	//assert.Equal(t, "config for redmine\n", actual.String())
}

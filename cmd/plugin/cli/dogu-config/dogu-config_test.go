package dogu_config

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

func TestCmd(t *testing.T) {
	originalK8sArgs := viper.GetViper().Get(util.CliTransportParamK8sArgs)
	defer viper.GetViper().Set(util.CliTransportParamK8sArgs, originalK8sArgs)
	t.Run("should print usage with namespace nil", func(t *testing.T) {
		// given
		viper.GetViper().Set(util.CliTransportParamK8sArgs, &genericclioptions.ConfigFlags{Namespace: nil})
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		sut := Cmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{})
		err := sut.Execute()

		// then
		assert.NoError(t, err)
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
	})
	t.Run("should print usage with namespace set", func(t *testing.T) {
		// given
		addressableNs := testNamespace
		viper.GetViper().Set(util.CliTransportParamK8sArgs, &genericclioptions.ConfigFlags{Namespace: &addressableNs})
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		sut := Cmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{})
		err := sut.Execute()

		// then
		assert.NoError(t, err)
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
	})
}

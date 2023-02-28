package dogu_config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

const testNamespace = "ecosystem"

func Test_defaultServiceFactory_create(t *testing.T) {
	t.Run("should fail to create rest config", func(t *testing.T) {
		// given
		cliConfigMock := newMockConfigTransporter(t)
		ns := testNamespace
		flags := &genericclioptions.ConfigFlags{
			Namespace: &ns,
		}
		cliConfigMock.EXPECT().Get(util.CliTransportParamK8sArgs).Return(flags).Once()
		sut := defaultServiceFactory{cliConfig: cliConfigMock}

		// when
		actual, err := sut.create(testDoguName)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "could not create rest config")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		cliConfigMock := newMockConfigTransporter(t)
		ns := testNamespace
		testKubeConfig := "./testdata/kubeConfig.valid"
		flags := &genericclioptions.ConfigFlags{
			Namespace:  &ns,
			KubeConfig: &testKubeConfig,
		}
		cliConfigMock.EXPECT().Get(util.CliTransportParamK8sArgs).Return(flags).Once()
		sut := defaultServiceFactory{cliConfig: cliConfigMock}

		// when
		actual, err := sut.create(testDoguName)

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewDoguConfigService(t *testing.T) {
	t.Run("should construct correctly", func(t *testing.T) {
		// given
		namespace := "test-namespace"

		// when
		actual, err := NewDoguConfigService(namespace, nil)

		// then
		require.NoError(t, err)

		assert.IsType(t, KubernetesPortForwarder{}, actual.portForwarder)
		actualPortForwarder := actual.portForwarder.(KubernetesPortForwarder)

		assert.NotEmpty(t, actualPortForwarder.LocalPort)

		expectedPortForwarder := KubernetesPortForwarder{
			RestConfig: nil,
			Type:       ServiceType,
			NamespacedName: types.NamespacedName{
				Namespace: "test-namespace",
				Name:      "etcd",
			},
			LocalPort:   actualPortForwarder.LocalPort,
			ClusterPort: 4001,
		}
		assert.Equal(t, expectedPortForwarder, actualPortForwarder)
	})
}

func TestDoguConfigService_Edit(t *testing.T) {
}

func TestDoguConfigService_Delete(t *testing.T) {
	t.Run("should", func(t *testing.T) {

	})
}

func TestDoguConfigService_getAllForDogu(t *testing.T) {
}

func TestDoguConfigService_GetValue(t *testing.T) {
}

func TestDoguConfigService_checkInstallStatus(t *testing.T) {
}

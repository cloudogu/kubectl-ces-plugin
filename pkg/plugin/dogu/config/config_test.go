package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	mocks2 "github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu/config/mocks"
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
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := mocks.NewPortForwarder(t)
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
		sut := DoguConfigService{
			registry:      nil,
			portForwarder: pfMock,
		}

		// when
		err := sut.Edit("official/some-dogu", "some-key", "some-value")

		// then
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail because not installed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(false, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		_ = sut.Edit("official/some-dogu", "some-key", "some-value")

		// then
		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
	})
	t.Run("should fail while setting key", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("Set", "some-key", "some-value").Once().Return(assert.AnError)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		_ = sut.Edit("official/some-dogu", "some-key", "some-value")

		// then
		assert.ErrorContains(t, err, "error while editing key 'some-key' for dogu 'official/some-dogu'")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("Set", "some-key", "some-value").Once().Return(nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		_ = sut.Edit("official/some-dogu", "some-key", "some-value")

		// then
		require.NoError(t, err)
	})
}

func TestDoguConfigService_Delete(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := mocks.NewPortForwarder(t)
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
		sut := DoguConfigService{
			registry:      nil,
			portForwarder: pfMock,
		}

		// when
		err := sut.Delete("official/some-dogu", "some-key")

		// then
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail because not installed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(false, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		_ = sut.Delete("official/some-dogu", "some-key")

		// then
		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
	})
	t.Run("should fail while deleting key", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("Delete", "some-key").Once().Return(assert.AnError)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		_ = sut.Delete("official/some-dogu", "some-key")

		// then
		assert.ErrorContains(t, err, "error while deleting key 'some-key' for dogu 'official/some-dogu'")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("Delete", "some-key").Once().Return(nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		_ = sut.Delete("official/some-dogu", "some-key")

		// then
		require.NoError(t, err)
	})
}

func TestDoguConfigService_getAllForDogu(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := mocks.NewPortForwarder(t)
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
		sut := DoguConfigService{
			registry:      nil,
			portForwarder: pfMock,
		}

		// when
		allKeys, err := sut.getAllForDogu("official/some-dogu")

		// then
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, allKeys)
	})
	t.Run("should fail because not installed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(false, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		allKeys, _ := sut.getAllForDogu("official/some-dogu")

		// then
		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
		assert.Nil(t, allKeys)
	})
	t.Run("should fail while getting all keys", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("GetAll").Once().Return(nil, assert.AnError)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		allKeys, _ := sut.getAllForDogu("official/some-dogu")

		// then
		assert.ErrorContains(t, err, "error while reading all keys for dogu 'official/some-dogu'")
		assert.Nil(t, allKeys)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		keys := map[string]string{
			"some-key":    "some-values",
			"another-key": "another-value",
		}
		configMock.On("GetAll").Once().Return(keys, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		allKeys, _ := sut.getAllForDogu("official/some-dogu")

		// then
		require.NoError(t, err)
		assert.Equal(t, keys, allKeys)
	})
}

func TestDoguConfigService_GetValue(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := mocks.NewPortForwarder(t)
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
		sut := DoguConfigService{
			registry:      nil,
			portForwarder: pfMock,
		}

		// when
		actual, err := sut.GetValue("official/some-dogu", "some-key")

		// then
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, actual)
	})
	t.Run("should fail because not installed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(false, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		actual, _ := sut.GetValue("official/some-dogu", "some-key")

		// then
		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
		assert.Empty(t, actual)
	})
	t.Run("should fail while deleting key", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("Get", "some-key").Once().Return("", assert.AnError)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		actual, _ := sut.GetValue("official/some-dogu", "some-key")

		// then
		assert.ErrorContains(t, err, "error while reading key 'some-key' for dogu 'official/some-dogu'")
		assert.Empty(t, actual)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		configMock := mocks2.NewConfigurationContext(t)
		configMock.On("Get", "some-key").Once().Return("some-value", nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock).
			On("DoguConfig", "official/some-dogu").Once().Return(configMock)
		pfMock := mocks.NewPortForwarder(t)
		var err error
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Run(func(args mock.Arguments) {
			fn, ok := args[0].(func() error)
			require.True(t, ok)
			err = fn()
		}).Return(nil)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		actual, _ := sut.GetValue("official/some-dogu", "some-key")

		// then
		require.NoError(t, err)
		assert.Equal(t, "some-value", actual)
	})
}

func TestDoguConfigService_checkInstallStatus(t *testing.T) {
	t.Run("should fail on check", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.
			On("IsEnabled", "official/some-dogu").Once().Return(false, assert.AnError)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.checkInstallStatus("official/some-dogu")

		// then
		require.Error(t, err)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot check if dogu 'official/some-dogu' is installed")
	})
	t.Run("should fail because not installed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.
			On("IsEnabled", "official/some-dogu").Once().Return(false, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.checkInstallStatus("official/some-dogu")

		// then
		require.Error(t, err)

		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.
			On("IsEnabled", "official/some-dogu").Once().Return(true, nil)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := mocks.NewPortForwarder(t)
		sut := DoguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.checkInstallStatus("official/some-dogu")

		// then
		require.NoError(t, err)
	})
}

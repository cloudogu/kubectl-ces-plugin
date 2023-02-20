package config

import (
	"github.com/cloudogu/kubectl-ces-plugin/pkg/portforward"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	mocks2 "github.com/cloudogu/cesapp-lib/registry/mocks"
)

func TestNewPortForwardedDoguConfigService(t *testing.T) {
	t.Run("should construct correctly", func(t *testing.T) {
		// given
		namespace := "test-namespace"

		// when
		actual, err := NewDoguConfigService(namespace, nil)
		assert.IsType(t, portforward.kubernetesPortForward{}, actual.portForwarder)

		// then
		require.NoError(t, err)

		assert.IsType(t, portforward.kubernetesPortForward{}, actual.portForwarder)
		actualPortForwarder := actual.portForwarder.(portforward.kubernetesPortForward)

		assert.NotEmpty(t, actualPortForwarder.localPort)

		expectedPortForwarder := portforward.kubernetesPortForward{
			restConfig: nil,
			typ:        portforward.ServiceType,
			name: types.NamespacedName{
				Namespace: "test-namespace",
				Name:      "etcd",
			},
			localPort:   actualPortForwarder.localPort,
			clusterPort: 4001,
		}
		assert.Equal(t, expectedPortForwarder, actualPortForwarder)
	})
}

func TestPortForwardedDoguConfigService_Edit(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{
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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.Edit("official/some-dogu", "some-key", "some-value")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.Edit("official/some-dogu", "some-key", "some-value")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.Edit("official/some-dogu", "some-key", "some-value")

		// then
		require.NoError(t, err)
	})
}

func TestPortForwardedDoguConfigService_Delete(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{
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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.Delete("official/some-dogu", "some-key")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.Delete("official/some-dogu", "some-key")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.Delete("official/some-dogu", "some-key")

		// then
		require.NoError(t, err)
	})
}

func TestPortForwardedDoguConfigService_GetAllForDogu(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := NewMockPortForwarder(t)
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
		sut := doguConfigService{
			registry:      nil,
			portForwarder: pfMock,
		}

		// when
		allKeys, err := sut.GetAllForDogu("official/some-dogu")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		allKeys, err := sut.GetAllForDogu("official/some-dogu")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		allKeys, err := sut.GetAllForDogu("official/some-dogu")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		allKeys, err := sut.GetAllForDogu("official/some-dogu")

		// then
		require.NoError(t, err)
		assert.Equal(t, keys, allKeys)
	})
}

func TestPortForwardedDoguConfigService_GetValue(t *testing.T) {
	t.Run("should fail during port-forward", func(t *testing.T) {
		// given
		pfMock := NewMockPortForwarder(t)
		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
		sut := doguConfigService{
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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		actual, err := sut.GetValue("official/some-dogu", "some-key")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		actual, err := sut.GetValue("official/some-dogu", "some-key")

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
		pfMock := NewMockPortForwarder(t)
		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
			return fn()
		}).Once()
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		actual, err := sut.GetValue("official/some-dogu", "some-key")

		// then
		require.NoError(t, err)
		assert.Equal(t, "some-value", actual)
	})
}

func TestPortForwardedDoguConfigService_checkInstallStatus(t *testing.T) {
	t.Run("should fail on check", func(t *testing.T) {
		// given
		doguRegistryMock := mocks2.NewDoguRegistry(t)
		doguRegistryMock.
			On("IsEnabled", "official/some-dogu").Once().Return(false, assert.AnError)
		registryMock := mocks2.NewRegistry(t)
		registryMock.
			On("DoguRegistry").Once().Return(doguRegistryMock)
		pfMock := NewMockPortForwarder(t)
		sut := doguConfigService{
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
		pfMock := NewMockPortForwarder(t)
		sut := doguConfigService{
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
		pfMock := NewMockPortForwarder(t)
		sut := doguConfigService{
			registry:      registryMock,
			portForwarder: pfMock,
		}

		// when
		err := sut.checkInstallStatus("official/some-dogu")

		// then
		require.NoError(t, err)
	})
}

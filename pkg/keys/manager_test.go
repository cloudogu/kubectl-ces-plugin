package keys

import (
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testDoguName    = "official/redmine"
	testKeyProvider = "pkcs1v15"
)

func TestNewKeyManager(t *testing.T) {
	t.Run("should fail to read key provider from registry", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("", assert.AnError)
		regMock := newMockCesRegistry(t)
		regMock.EXPECT().GlobalConfig().Return(globalConfigMock)

		// when
		_, err := NewKeyManager(regMock, testDoguName)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not create etcd key manager: could not read key type")
	})
	t.Run("should fail to create key provider", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("banane", nil)
		regMock := newMockCesRegistry(t)
		regMock.EXPECT().GlobalConfig().Return(globalConfigMock)

		// when
		_, err := NewKeyManager(regMock, testDoguName)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not create key provider")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("key_provider").Return(testKeyProvider, nil)
		doguConfigMock := newMockConfigurationContext(t)
		regMock := newMockCesRegistry(t)
		regMock.EXPECT().GlobalConfig().Return(globalConfigMock)
		regMock.EXPECT().DoguConfig(testDoguName).Return(doguConfigMock)

		// when
		actual, err := NewKeyManager(regMock, testDoguName)

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func Test_etcdKeyManager_ExistsPublicKey(t *testing.T) {
	t.Run("should fail on checking if public key exists", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(false, assert.AnError)
		sut := etcdKeyManager{doguConfig: doguConfigMock}

		// when
		_, err := sut.ExistsPublicKey()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check if public key exists")
	})
	t.Run("should succeed if public key does not exist", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(false, nil)
		sut := etcdKeyManager{doguConfig: doguConfigMock}

		// when
		exists, err := sut.ExistsPublicKey()

		// then
		require.NoError(t, err)
		assert.Equal(t, false, exists)
	})
	t.Run("should succeed if public key exists", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(true, nil)
		sut := etcdKeyManager{doguConfig: doguConfigMock}

		// when
		exists, err := sut.ExistsPublicKey()

		// then
		require.NoError(t, err)
		assert.Equal(t, true, exists)
	})
}

func Test_etcdKeyManager_GetPublicKey(t *testing.T) {
	t.Run("should fail on checking if public key exists", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(false, assert.AnError)
		sut := etcdKeyManager{doguConfig: doguConfigMock}

		// when
		_, err := sut.GetPublicKey()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check if public key exists")
	})
	t.Run("should fail if public key does not exist", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(false, nil)
		sut := etcdKeyManager{doguConfig: doguConfigMock}

		// when
		_, err := sut.GetPublicKey()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not find public key in configuration context")
	})
	t.Run("should fail on fetching public key", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(true, nil)
		doguConfigMock.EXPECT().Get(registry.KeyDoguPublicKey).Return("", assert.AnError)
		sut := etcdKeyManager{doguConfig: doguConfigMock}

		// when
		_, err := sut.GetPublicKey()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to fetch public key from registry")
	})
	t.Run("should fail on reading key", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(true, nil)
		doguConfigMock.EXPECT().Get(registry.KeyDoguPublicKey).Return("---PUBLIC KEY---", nil)
		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().ReadPublicKeyFromString("---PUBLIC KEY---").Return(nil, assert.AnError)
		sut := etcdKeyManager{
			doguConfig:  doguConfigMock,
			keyProvider: keyProviderMock,
		}

		// when
		_, err := sut.GetPublicKey()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create public key from string")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		doguConfigMock := newMockConfigurationContext(t)
		doguConfigMock.EXPECT().Exists(registry.KeyDoguPublicKey).Return(true, nil)
		doguConfigMock.EXPECT().Get(registry.KeyDoguPublicKey).Return("---PUBLIC KEY---", nil)
		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().ReadPublicKeyFromString("---PUBLIC KEY---").Return(&keys.PublicKey{}, nil)
		sut := etcdKeyManager{
			doguConfig:  doguConfigMock,
			keyProvider: keyProviderMock,
		}

		// when
		_, err := sut.GetPublicKey()

		// then
		require.NoError(t, err)
	})
}

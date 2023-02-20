package config

const testNameSpace = "test-namespace"
const testDoguName = "official/ldap"

//
// func TestPortForwardedDoguConfigService_Edit(t *testing.T) {
// 	t.Run("should fail during port-forward", func(t *testing.T) {
// 		// given
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).Return(assert.AnError).Once()
// 		sut := doguConfigService{
// 			registry:      nil,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Edit(testDoguName, "some-key", "some-value")
//
// 		// then
// 		assert.ErrorIs(t, err, assert.AnError)
// 	})
// 	t.Run("should fail because not installed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(false, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Edit(testDoguName, "some-key", "some-value")
//
// 		// then
// 		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
// 	})
// 	t.Run("should fail while setting key", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("Set", "some-key", "some-value").Once().Return(assert.AnError)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Edit(testDoguName, "some-key", "some-value")
//
// 		// then
// 		assert.ErrorContains(t, err, "error while editing key 'some-key' for dogu 'official/some-dogu'")
// 	})
// 	t.Run("should succeed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("Set", "some-key", "some-value").Once().Return(nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Edit(testDoguName, "some-key", "some-value")
//
// 		// then
// 		require.NoError(t, err)
// 	})
// }
//
// func TestPortForwardedDoguConfigService_Delete(t *testing.T) {
// 	t.Run("should fail during port-forward", func(t *testing.T) {
// 		// given
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).Return(assert.AnError).Once()
// 		sut := doguConfigService{
// 			registry:      nil,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Delete(testDoguName, "some-key")
//
// 		// then
// 		assert.ErrorIs(t, err, assert.AnError)
// 	})
// 	t.Run("should fail because not installed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(false, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Delete(testDoguName, "some-key")
//
// 		// then
// 		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
// 	})
// 	t.Run("should fail while deleting key", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("Delete", "some-key").Once().Return(assert.AnError)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Delete(testDoguName, "some-key")
//
// 		// then
// 		assert.ErrorContains(t, err, "error while deleting key 'some-key' for dogu 'official/some-dogu'")
// 	})
// 	t.Run("should succeed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("Delete", "some-key").Once().Return(nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.Delete(testDoguName, "some-key")
//
// 		// then
// 		require.NoError(t, err)
// 	})
// }
//
// func TestPortForwardedDoguConfigService_GetAllForDogu(t *testing.T) {
// 	t.Run("should fail during port-forward", func(t *testing.T) {
// 		// given
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
// 		sut := doguConfigService{
// 			registry:      nil,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		allKeys, err := sut.GetAllForDogu(testDoguName)
//
// 		// then
// 		assert.ErrorIs(t, err, assert.AnError)
// 		assert.Nil(t, allKeys)
// 	})
// 	t.Run("should fail because not installed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(false, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		allKeys, err := sut.GetAllForDogu(testDoguName)
//
// 		// then
// 		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
// 		assert.Nil(t, allKeys)
// 	})
// 	t.Run("should fail while getting all keys", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("GetAll").Once().Return(nil, assert.AnError)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		allKeys, err := sut.GetAllForDogu(testDoguName)
//
// 		// then
// 		assert.ErrorContains(t, err, "error while reading all keys for dogu 'official/some-dogu'")
// 		assert.Nil(t, allKeys)
// 	})
// 	t.Run("should succeed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		keys := map[string]string{
// 			"some-key":    "some-values",
// 			"another-key": "another-value",
// 		}
// 		configMock.On("GetAll").Once().Return(keys, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		allKeys, err := sut.GetAllForDogu(testDoguName)
//
// 		// then
// 		require.NoError(t, err)
// 		assert.Equal(t, keys, allKeys)
// 	})
// }
//
// func TestPortForwardedDoguConfigService_GetValue(t *testing.T) {
// 	t.Run("should fail during port-forward", func(t *testing.T) {
// 		// given
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.On("ExecuteWithPortForward", mock.Anything).Once().Return(assert.AnError)
// 		sut := doguConfigService{
// 			registry:      nil,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		actual, err := sut.GetValue(testDoguName, "some-key")
//
// 		// then
// 		assert.ErrorIs(t, err, assert.AnError)
// 		assert.Empty(t, actual)
// 	})
// 	t.Run("should fail because not installed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(false, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		actual, err := sut.GetValue(testDoguName, "some-key")
//
// 		// then
// 		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
// 		assert.Empty(t, actual)
// 	})
// 	t.Run("should fail while deleting key", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("Get", "some-key").Once().Return("", assert.AnError)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		actual, err := sut.GetValue(testDoguName, "some-key")
//
// 		// then
// 		assert.ErrorContains(t, err, "error while reading key 'some-key' for dogu 'official/some-dogu'")
// 		assert.Empty(t, actual)
// 	})
// 	t.Run("should succeed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		configMock := mocks2.NewConfigurationContext(t)
// 		configMock.On("Get", "some-key").Once().Return("some-value", nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock).
// 			On("DoguConfig", testDoguName).Once().Return(configMock)
// 		pfMock := NewMockPortForwarder(t)
// 		pfMock.EXPECT().ExecuteWithPortForward(mock.Anything).RunAndReturn(func(fn func() error) error {
// 			return fn()
// 		}).Once()
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		actual, err := sut.GetValue(testDoguName, "some-key")
//
// 		// then
// 		require.NoError(t, err)
// 		assert.Equal(t, "some-value", actual)
// 	})
// }
//
// func TestPortForwardedDoguConfigService_checkInstallStatus(t *testing.T) {
// 	t.Run("should fail on check", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.
// 			On("IsEnabled", testDoguName).Once().Return(false, assert.AnError)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.checkInstallStatus(testDoguName)
//
// 		// then
// 		require.Error(t, err)
//
// 		assert.ErrorIs(t, err, assert.AnError)
// 		assert.ErrorContains(t, err, "cannot check if dogu 'official/some-dogu' is installed")
// 	})
// 	t.Run("should fail because not installed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.
// 			On("IsEnabled", testDoguName).Once().Return(false, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.checkInstallStatus(testDoguName)
//
// 		// then
// 		require.Error(t, err)
//
// 		assert.ErrorContains(t, err, "dogu 'official/some-dogu' is not installed")
// 	})
// 	t.Run("should succeed", func(t *testing.T) {
// 		// given
// 		doguRegistryMock := mocks2.NewDoguRegistry(t)
// 		doguRegistryMock.
// 			On("IsEnabled", testDoguName).Once().Return(true, nil)
// 		registryMock := mocks2.NewRegistry(t)
// 		registryMock.
// 			On("DoguRegistry").Once().Return(doguRegistryMock)
// 		pfMock := NewMockPortForwarder(t)
// 		sut := doguConfigService{
// 			registry:      registryMock,
// 			portForwarder: pfMock,
// 		}
//
// 		// when
// 		err := sut.checkInstallStatus(testDoguName)
//
// 		// then
// 		require.NoError(t, err)
// 	})
// }

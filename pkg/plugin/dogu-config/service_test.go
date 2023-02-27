package dogu_config

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func Test_New(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		// given

		// when
		_, err := New(testDoguName, testNameSpace, &rest.Config{})

		// then
		require.NoError(t, err)
	})
}

func Test_doguConfigService_Edit(t *testing.T) {
	t.Run("should fail during delegate", func(t *testing.T) {
		// given
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Edit("test-key", false)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail with no matching configuration fields", func(t *testing.T) {
		// given
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{{
				Name: "nonmatching",
			}},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Edit("test-key", false)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu 'ldap' has no matching configuration fields for key 'test-key'")
	})
	t.Run("should fail on edit configuration of matching field", func(t *testing.T) {
		// given
		configFields := []core.ConfigurationField{{
			Name: "matching",
		}}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: configFields,
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().EditConfiguration(configFields, false).Return(assert.AnError).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Edit("match", false)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should succeed on edit configuration of matching field", func(t *testing.T) {
		// given
		configFields := []core.ConfigurationField{{
			Name: "matching",
		}}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: configFields,
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().EditConfiguration(configFields, false).Return(nil).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Edit("match", false)

		// then
		require.NoError(t, err)
	})
}

func Test_matchConfigurationFields(t *testing.T) {
	t.Run("should return all config fields if key is empty", func(t *testing.T) {
		// given
		configFields := []core.ConfigurationField{
			{Name: "some-field"},
			{Name: "another-field"},
		}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: configFields,
		}

		// when
		actual := matchConfigurationFields(dogu, "")

		// then
		assert.Equal(t, configFields, actual)
	})
	t.Run("should return matching config fields", func(t *testing.T) {
		// given
		configFields := []core.ConfigurationField{
			{Name: "matching-path/some-field"},
			{Name: "matching-path/another-field"},
		}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: append(configFields, core.ConfigurationField{Name: "non-matching-path/another-field"}),
		}

		// when
		actual := matchConfigurationFields(dogu, "matching")

		// then
		assert.Equal(t, configFields, actual)
	})
}

func Test_doguConfigService_Set(t *testing.T) {
	t.Run("should fail on delegate", func(t *testing.T) {
		// given
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Set("test-key", "test-value")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail with configuration field not found", func(t *testing.T) {
		// given
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{{
				Name: "nonmatching",
			}},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Set("test-key", "test-value")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu 'ldap' has no configuration field for key 'test-key'")
	})
	t.Run("should fail on set value", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: []core.ConfigurationField{testField},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().SetFieldToValue(testField, "test-value", false).Return(assert.AnError).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Set("test-key", "test-value")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should succeed on set value", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: []core.ConfigurationField{testField},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().SetFieldToValue(testField, "test-value", false).Return(nil).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Set("test-key", "test-value")

		// then
		require.NoError(t, err)
	})
}

func Test_doguConfigService_Delete(t *testing.T) {
	t.Run("should fail on delegate", func(t *testing.T) {
		// given
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Delete("test-key")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail with configuration field not found", func(t *testing.T) {
		// given
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{{
				Name: "nonmatching",
			}},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Delete("test-key")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu 'ldap' has no configuration field for key 'test-key'")
	})
	t.Run("should fail to delete key", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: []core.ConfigurationField{testField},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().DeleteField(testField).Return(assert.AnError).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Delete("test-key")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should succeed to delete key", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: []core.ConfigurationField{testField},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().DeleteField(testField).Return(nil).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Delete("test-key")

		// then
		require.NoError(t, err)
	})
}

func Test_doguConfigService_List(t *testing.T) {
	t.Run("should fail on delegate", func(t *testing.T) {
		// given
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		_, err := sut.List()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail to list config fields", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{
				testField,
				{Name: "some-key"},
			},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().GetCurrentValue(testField).Return("", assert.AnError).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		_, err := sut.List()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should succeed to list config fields", func(t *testing.T) {
		// given
		testField1 := core.ConfigurationField{Name: "some-key"}
		testField2 := core.ConfigurationField{Name: "another-key"}
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{
				testField1,
				testField2,
			},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().GetCurrentValue(testField1).Return("some-value", nil).Once()
		mockEditor.EXPECT().GetCurrentValue(testField2).Return("another-value", nil).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}
		expected := map[string]string{
			"some-key":    "some-value",
			"another-key": "another-value",
		}

		// when
		actual, err := sut.List()

		// then
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_doguConfigService_GetValue(t *testing.T) {
	t.Run("should fail on delegate", func(t *testing.T) {
		// given
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).Return(assert.AnError).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		_, err := sut.GetValue("test-key")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail with configuration field not found", func(t *testing.T) {
		// given
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{{
				Name: "nonmatching",
			}},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		_, err := sut.GetValue("test-key")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu 'ldap' has no configuration field for key 'test-key'")
	})
	t.Run("should fail to get value for key", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: []core.ConfigurationField{testField},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().GetCurrentValue(testField).Return("", assert.AnError).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		_, err := sut.GetValue("test-key")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should succeed to get value for key", func(t *testing.T) {
		// given
		testField := core.ConfigurationField{Name: "test-key"}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: []core.ConfigurationField{testField},
		}
		mockEditor := newMockDoguConfigurationEditor(t)
		mockEditor.EXPECT().GetCurrentValue(testField).Return("test-value", nil).Once()
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).RunAndReturn(func(payload func(*core.Dogu, doguConfigurationEditor) error) error {
			return payload(dogu, mockEditor)
		}).Once()
		sut := doguConfigService{delegator: delegator}

		// when
		actual, err := sut.GetValue("test-key")

		// then
		require.NoError(t, err)
		assert.Equal(t, "test-value", actual)
	})
}

func Test_configurationFieldByKey(t *testing.T) {
	t.Run("should not find config field", func(t *testing.T) {
		// given
		configFields := []core.ConfigurationField{
			{Name: "some-field"},
			{Name: "another-field"},
		}
		dogu := &core.Dogu{
			Name:          testDoguName,
			Configuration: configFields,
		}

		// when
		found, actual := configurationFieldByKey(dogu, "no-match")

		// then
		assert.Equal(t, found, false)
		assert.Nil(t, actual)
	})
	t.Run("should find config field", func(t *testing.T) {
		// given
		matchingField := core.ConfigurationField{Name: "matching-field"}
		dogu := &core.Dogu{
			Name: testDoguName,
			Configuration: []core.ConfigurationField{
				{Name: "some-field"},
				{Name: "another-field"},
				matchingField,
			},
		}

		// when
		found, actual := configurationFieldByKey(dogu, "matching-field")

		// then
		assert.Equal(t, found, true)
		assert.Equal(t, matchingField, *actual)
	})
}

package dogu_config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_deleteCmd(t *testing.T) {
	t.Run("should get config value", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().Delete(testConfigKey).Return(nil).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := deleteCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey})
		err := sut.Execute()

		// then
		require.NoError(t, err, "command should be successful")
		assert.Empty(t, outBuf.String())
		assert.Empty(t, errBuf.String())
	})

	t.Run("should return error from configService", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().Delete(testConfigKey).Return(assert.AnError).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := deleteCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey})
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.ErrorContains(t, err, fmt.Sprintf("cannot delete config key '%s' in delete dogu config command", testConfigKey))
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error that the config service cannot be created", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(nil, assert.AnError).Once()

		sut := deleteCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey})
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.ErrorContains(t, err, "cannot create config service")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should fail with too few Arguments", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		serviceFactoryMock := newMockServiceFactory(t)
		sut := deleteCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.EqualError(t, err, "accepts 2 arg(s), received 0")
	})
}

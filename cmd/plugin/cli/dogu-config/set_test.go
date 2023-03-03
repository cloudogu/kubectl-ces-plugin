package dogu_config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setCmd(t *testing.T) {
	t.Run("should set config value", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().Set(testConfigKey, testConfigValue).Return(nil).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := setCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey, testConfigValue})
		err := sut.Execute()

		// then
		assert.NoError(t, err, "command should be successful")
		assert.Empty(t, outBuf.String())
		assert.Empty(t, errBuf.String())
	})

	t.Run("should return error from configService", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().Set(testConfigKey, testConfigValue).Return(assert.AnError).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := setCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey, testConfigValue})
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.ErrorContains(t, err, fmt.Sprintf("cannot set config key '%s'", testConfigKey))
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error that the config service cannot be created", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(nil, assert.AnError).Once()

		sut := setCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey, testConfigValue})
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
		sut := setCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.EqualError(t, err, "accepts 3 arg(s), received 0")
	})
}

package dogu_config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDoguName  = "redmine"
	testConfigKey = "redmineKey"
)

func Test_editCmd(t *testing.T) {
	t.Run("should edit all config values of dogu", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().Edit("", false).Return(nil).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := editCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName})
		err := sut.Execute()

		// then
		assert.NoError(t, err, "command should be successful")
		assert.Empty(t, outBuf.String())
		assert.Empty(t, errBuf.String())
	})

	t.Run("should edit specific config value", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().Edit(testConfigKey, false).Return(nil).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := editCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey})
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
		serviceMock.EXPECT().Edit(testConfigKey, false).Return(assert.AnError).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := editCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName, testConfigKey})
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.ErrorContains(t, err, "cannot edit config keys")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error that the config service cannot be created", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(nil, assert.AnError).Once()

		sut := editCmd(serviceFactoryMock)
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
		sut := editCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.EqualError(t, err, "accepts between 1 and 2 arg(s), received 0")
	})
}

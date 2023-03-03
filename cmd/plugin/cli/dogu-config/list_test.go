package dogu_config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfigList = map[string]string{
	"testKey1": "testValue1",
	"testKey2": "testValue2",
}

func Test_listCmd(t *testing.T) {
	t.Run("should list all config keys", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().List().Return(testConfigList, nil).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := listCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName})
		err := sut.Execute()

		// then
		assert.NoError(t, err, "command should be successful")
		assert.Contains(t, outBuf.String(), "testKey1: testValue1\n")
		assert.Contains(t, outBuf.String(), "testKey2: testValue2\n")
		assert.Empty(t, errBuf.String())
	})

	t.Run("should return error from configService", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceMock := newMockDoguConfigService(t)
		serviceMock.EXPECT().List().Return(nil, assert.AnError).Once()
		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(serviceMock, nil).Once()

		sut := listCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName})
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.ErrorContains(t, err, "cannot list config in list dogu config command")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error that the config service cannot be created", func(t *testing.T) {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		serviceFactoryMock := newMockServiceFactory(t)
		serviceFactoryMock.EXPECT().create(testDoguName).Return(nil, assert.AnError).Once()

		sut := listCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{testDoguName})
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
		sut := listCmd(serviceFactoryMock)
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		assert.Contains(t, outBuf.String(), "Usage:", "should have usage output")
		assert.Contains(t, errBuf.String(), err.Error(), "should contain error output")
		assert.EqualError(t, err, "accepts 1 arg(s), received 0")
	})
}

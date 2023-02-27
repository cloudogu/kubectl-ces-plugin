package dogu_config

import (
	"bytes"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func (s *DoguConfigCLITestSuite) Test_setCmd() {
	const doguName = "redmine"
	const configKey = "redmineKey"
	const configValue = "redmineValue"
	s.Run("should set config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := setCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		doguConfigServiceFactoryMock.EXPECT().Set(configKey, configValue).Return(nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		sut.SetArgs([]string{doguName, configKey, configValue})
		err := sut.Execute()

		// then
		s.NoError(err, "command should be successful")
		s.Empty(outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := setCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		doguConfigServiceFactoryMock.EXPECT().Set(configKey, configValue).Return(assert.AnError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		sut.SetArgs([]string{doguName, configKey, configValue})
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, fmt.Sprintf("cannot set config key '%s'", configKey))
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := setCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		doguConfigServiceFactory = errorDoguConfigServiceFactory(assert.AnError)

		// when
		sut.SetArgs([]string{doguName, configKey, configValue})
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, "cannot create config service")
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should fail with too few Arguments", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := setCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 3 arg(s), received 0")
	})
}

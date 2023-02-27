package dogu_config

import (
	"bytes"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func (s *DoguConfigCLITestSuite) Test_getCmd() {
	s.Run("should get config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := getCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		configValue := "redmineValue"
		doguConfigServiceFactoryMock.EXPECT().GetValue(configKey).Return(configValue, nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		sut.SetArgs([]string{doguName, configKey})
		err := sut.Execute()

		// then
		s.NoError(err, "command should be successful")
		s.Equal(configValue, outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := getCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		doguConfigServiceFactoryMock.EXPECT().GetValue(configKey).Return("", assert.AnError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		sut.SetArgs([]string{doguName, configKey})
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, fmt.Sprintf("cannot get config key '%s' in get dogu config command", configKey))
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := getCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)
		doguName := "redmine"

		configKey := "redmineKey"
		doguConfigServiceFactory = errorDoguConfigServiceFactory(assert.AnError)

		// when
		sut.SetArgs([]string{doguName, configKey})
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
		sut := getCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 2 arg(s), received 0")
	})
}

package dogu_config

import (
	"bytes"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func (s *DoguConfigCLITestSuite) Test_deleteCmd() {
	s.Run("should get config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := deleteCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		doguConfigServiceFactoryMock.EXPECT().Delete(configKey).Return(nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{doguName, configKey})
		err := configCmd.Execute()

		// then
		s.NoError(err, "command should be successful")
		s.Empty(outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := deleteCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		doguConfigServiceFactoryMock.EXPECT().Delete(configKey).Return(assert.AnError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{doguName, configKey})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, fmt.Sprintf("cannot delete config key '%s' in delete dogu config command", configKey))
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := deleteCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"

		configKey := "redmineKey"
		doguConfigServiceFactory = errorDoguConfigServiceFactory(assert.AnError)

		// when
		configCmd.SetArgs([]string{doguName, configKey})
		err := configCmd.Execute()

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
		configCmd := deleteCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		// when
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 2 arg(s), received 0")
	})
}

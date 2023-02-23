package dogu_config

import (
	"bytes"

	"github.com/stretchr/testify/assert"
)

func (s *DoguConfigCLITestSuite) Test_listCmd() {
	s.Run("should list all config keys", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := listCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		returnedConfig := make(map[string]string)
		returnedConfig["testKey1"] = "testValue1"
		returnedConfig["testKey2"] = "testValue2"
		doguConfigServiceFactoryMock.EXPECT().List().Return(returnedConfig, nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{doguName})
		err := configCmd.Execute()

		// then
		s.NoError(err, "command should be successful")
		s.Contains(outBuf.String(), "testKey1: testValue1\n")
		s.Contains(outBuf.String(), "testKey2: testValue2\n")
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := listCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		doguConfigServiceFactoryMock.EXPECT().List().Return(nil, assert.AnError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{doguName})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, "cannot list config in list dogu config command")
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := listCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactory = errorDoguConfigServiceFactory(assert.AnError)

		// when
		configCmd.SetArgs([]string{doguName})
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
		configCmd := listCmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		// when
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 1 arg(s), received 0")
	})
}

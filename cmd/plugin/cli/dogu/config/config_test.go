package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DoguConfigCLITestSuite struct {
	suite.Suite
	originalFactory func(viper *viper.Viper) (DoguConfigService, error)
}

func TestDoguConfigCLITestSuite(t *testing.T) {
	suite.Run(t, new(DoguConfigCLITestSuite))
}

func (s *DoguConfigCLITestSuite) SetupSuite() {
	s.originalFactory = DoguConfigServiceFactory
}

func (s *DoguConfigCLITestSuite) TearDownSuite() {
	DoguConfigServiceFactory = s.originalFactory
}

func (s *DoguConfigCLITestSuite) Test_getAllForDoguCmd() {
	s.Run("should list all config keys", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		returnedConfig := make(map[string]string)
		returnedConfig["testKey1"] = "testValue1"
		returnedConfig["testKey2"] = "testValue2"
		mock.EXPECT().GetAllForDogu(doguName).Return(returnedConfig, nil).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"list", doguName})
		err := configCmd.Execute()

		//then
		s.NoError(err, "command should be successful")
		s.Contains(outBuf.String(), "testKey1: testValue1\n")
		s.Contains(outBuf.String(), "testKey2: testValue2\n")
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		expectedError := errors.New("configService error")
		mock.EXPECT().GetAllForDogu(doguName).Return(nil, expectedError).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"list", doguName})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config list <doguName> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot list config in list dogu config command: configService error")
	})

	s.Run("should return error that the config service cannot be created", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		expectedError := errors.New("create configService error")
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return nil, expectedError
		}

		//when
		configCmd.SetArgs([]string{"list", doguName})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config list <doguName> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in list dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		//when
		configCmd.SetArgs([]string{"list"})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config list <doguName> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 1 arg(s), received 0")
	})
}

func (s *DoguConfigCLITestSuite) Test_getCmd() {
	s.Run("should get config value", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		configKey := "redmineKey"
		configValue := "redmineValue"
		mock.EXPECT().GetValue(doguName, configKey).Return(configValue, nil).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"get", doguName, configKey})
		err := configCmd.Execute()

		//then
		s.NoError(err, "command should be successful")
		s.Equal(configValue, outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		configKey := "redmineKey"
		expectedError := errors.New("configService error")
		mock.EXPECT().GetValue(doguName, configKey).Return("", expectedError).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"get", doguName, configKey})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config get <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, fmt.Sprintf("cannot get config key '%s' in get dogu config command: configService error", configKey))
	})

	s.Run("should return error that the config service cannot be created", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		configKey := "redmineKey"
		expectedError := errors.New("create configService error")
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return nil, expectedError
		}

		//when
		configCmd.SetArgs([]string{"get", doguName, configKey})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config get <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in get dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		//when
		configCmd.SetArgs([]string{"get"})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config get <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 2 arg(s), received 0")
	})
}

func (s *DoguConfigCLITestSuite) Test_editCmd() {
	s.Run("should set config value", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		configKey := "redmineKey"
		configValue := "redmineValue"
		mock.EXPECT().Edit(doguName, configKey, configValue).Return(nil).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"edit", doguName, configKey, configValue})
		err := configCmd.Execute()

		//then
		s.NoError(err, "command should be successful")
		s.Empty(outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		configKey := "redmineKey"
		configValue := "redmineValue"
		expectedError := errors.New("configService error")
		mock.EXPECT().Edit(doguName, configKey, configValue).Return(expectedError).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"edit", doguName, configKey, configValue})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(),
			"Usage:\n  config edit <doguName> <configKey> <configValue> [flags]",
			"should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, fmt.Sprintf("cannot set config key '%s' in edit dogu config command: configService error", configKey))
	})

	s.Run("should return error that the config service cannot be created", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		configKey := "redmineKey"
		configValue := "redmineValue"
		expectedError := errors.New("create configService error")
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return nil, expectedError
		}

		//when
		configCmd.SetArgs([]string{"edit", doguName, configKey, configValue})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(),
			"Usage:\n  config edit <doguName> <configKey> <configValue> [flags]",
			"should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in get dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		//when
		configCmd.SetArgs([]string{"edit"})
		err := configCmd.Execute()

		//then
		s.Contains(
			outBuf.String(), "Usage:\n  config edit <doguName> <configKey> <configValue> [flags]",
			"should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 3 arg(s), received 0")
	})
}

func (s *DoguConfigCLITestSuite) Test_deleteCmd() {
	s.Run("should get config value", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		configKey := "redmineKey"
		mock.EXPECT().Delete(doguName, configKey).Return(nil).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"delete", doguName, configKey})
		err := configCmd.Execute()

		//then
		s.NoError(err, "command should be successful")
		s.Empty(outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		mock := NewMockDoguConfigService(s.T())
		configKey := "redmineKey"
		expectedError := errors.New("configService error")
		mock.EXPECT().Delete(doguName, configKey).Return(expectedError).Once()
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return mock, nil
		}

		//when
		configCmd.SetArgs([]string{"delete", doguName, configKey})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config delete <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, fmt.Sprintf("cannot delete config key '%s' in delete dogu config command: configService error", configKey))
	})

	s.Run("should return error that the config service cannot be created", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set("doguName", doguName)

		configKey := "redmineKey"
		expectedError := errors.New("create configService error")
		DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
			return nil, expectedError
		}

		//when
		configCmd.SetArgs([]string{"delete", doguName, configKey})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config delete <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in get dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		//given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		//when
		configCmd.SetArgs([]string{"delete"})
		err := configCmd.Execute()

		//then
		s.Contains(outBuf.String(), "Usage:\n  config delete <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 2 arg(s), received 0")
	})
}
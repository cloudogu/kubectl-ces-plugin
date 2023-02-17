package config

import (
	"bytes"
	"errors"
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
		s.Equal("", errBuf.String())
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

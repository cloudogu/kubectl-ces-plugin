package config

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CliTestSuite struct {
	suite.Suite
	originalFactory func(viper *viper.Viper) (DoguConfigService, error)
}

func TestCliTestSuite(t *testing.T) {
	suite.Run(t, new(CliTestSuite))
}

func (s *CliTestSuite) SetupSuite() {
	s.originalFactory = DoguConfigServiceFactory
}

func (s *CliTestSuite) TearDownSuite() {
	DoguConfigServiceFactory = s.originalFactory
}

func (s *CliTestSuite) Test_getAllForDoguCmd() {
	//given
	actual := new(bytes.Buffer)
	configCmd := Cmd()
	configCmd.SetOut(actual)
	configCmd.SetErr(actual)
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
	s.Contains(actual.String(), "testKey1: testValue1\n")
	s.Contains(actual.String(), "testKey2: testValue2\n")
}

package dogu_config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DoguConfigCLITestSuite struct {
	suite.Suite
	originalFactory func(doguName string) (doguConfigService, error)
}

func TestDoguConfigCLITestSuite(t *testing.T) {
	suite.Run(t, new(DoguConfigCLITestSuite))
}

func (s *DoguConfigCLITestSuite) SetupSuite() {
	s.originalFactory = doguConfigServiceFactory
}

func (s *DoguConfigCLITestSuite) TearDownSuite() {
	doguConfigServiceFactory = s.originalFactory
}

func noopDoguConfigServiceFactory(configServiceMock *mockDoguConfigService) func(doguName string) (doguConfigService, error) {
	return func(doguName string) (doguConfigService, error) {
		return configServiceMock, nil
	}
}

func errorDoguConfigServiceFactory(expectedError error) func(doguName string) (doguConfigService, error) {
	return func(doguName string) (doguConfigService, error) {
		return nil, expectedError
	}
}

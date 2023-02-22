package config

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

type DoguConfigCLITestSuite struct {
	suite.Suite
	originalFactory func(doguName, k8sNamespace string, restConfig *rest.Config) (doguConfigService, error)
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

func (s *DoguConfigCLITestSuite) Test_getAllForDoguCmd() {
	s.Run("should list all config keys", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		returnedConfig := make(map[string]string)
		returnedConfig["testKey1"] = "testValue1"
		returnedConfig["testKey2"] = "testValue2"
		doguConfigServiceFactoryMock.EXPECT().List().Return(returnedConfig, nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"list", doguName})
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
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		expectedError := errors.New("configService error")
		doguConfigServiceFactoryMock.EXPECT().List().Return(nil, expectedError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"list", doguName})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config list <doguName> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot list config in list dogu config command: configService error")
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		expectedError := errors.New("create configService error")
		doguConfigServiceFactory = errorDoguConfigServiceFactory(expectedError)

		// when
		configCmd.SetArgs([]string{"list", doguName})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config list <doguName> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in list dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		// when
		configCmd.SetArgs([]string{"list"})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config list <doguName> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 1 arg(s), received 0")
	})
}

func (s *DoguConfigCLITestSuite) Test_getCmd() {
	s.Run("should get config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		configValue := "redmineValue"
		doguConfigServiceFactoryMock.EXPECT().GetValue(configKey).Return(configValue, nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"get", doguName, configKey})
		err := configCmd.Execute()

		// then
		s.NoError(err, "command should be successful")
		s.Equal(configValue, outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		expectedError := errors.New("configService error")
		doguConfigServiceFactoryMock.EXPECT().GetValue(configKey).Return("", expectedError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"get", doguName, configKey})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config get <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, fmt.Sprintf("cannot get config key '%s' in get dogu config command: configService error", configKey))
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		configKey := "redmineKey"
		expectedError := errors.New("create configService error")
		doguConfigServiceFactory = errorDoguConfigServiceFactory(expectedError)

		// when
		configCmd.SetArgs([]string{"get", doguName, configKey})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config get <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in get dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		// when
		configCmd.SetArgs([]string{"get"})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config get <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 2 arg(s), received 0")
	})
}

func (s *DoguConfigCLITestSuite) Test_editCmd() {
	s.Run("should set config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		const configKey = "redmineKey"
		const deleteOnEdit = false
		doguConfigServiceFactoryMock.EXPECT().Edit(configKey, deleteOnEdit).Return(nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"edit", doguName, configKey})
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
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		const configKey = "redmineKey"
		expectedError := errors.New("configService error")
		const deleteOnEdit = false
		doguConfigServiceFactoryMock.EXPECT().Edit(configKey, deleteOnEdit).Return(expectedError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"edit", doguName, configKey})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(),
			"Usage:\n  config edit <doguName> <configKey> <configValue> [flags]",
			"should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, fmt.Sprintf("cannot set config key '%s' in edit dogu config command: configService error", configKey))
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		configKey := "redmineKey"
		configValue := "redmineValue"
		expectedError := errors.New("create configService error")
		doguConfigServiceFactory = errorDoguConfigServiceFactory(expectedError)

		// when
		configCmd.SetArgs([]string{"edit", doguName, configKey, configValue})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(),
			"Usage:\n  config edit <doguName> <configKey> <configValue> [flags]",
			"should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in get dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		viperInjectKubeEnvironment()

		// when
		configCmd.SetArgs([]string{"edit"})
		err := configCmd.Execute()

		// then
		s.Contains(
			outBuf.String(), "Usage:\n  config edit <doguName> <configKey> <configValue> [flags]",
			"should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 3 arg(s), received 0")
	})
}

func viperInjectKubeEnvironment() {
	k8sFlags := &genericclioptions.ConfigFlags{}
	validKubeConfig := "./testdata/kubeConfig.valid"
	k8sFlags.KubeConfig = &validKubeConfig
	viper.GetViper().Set(util.CliTransportParamK8sArgs, k8sFlags)
}

func (s *DoguConfigCLITestSuite) Test_deleteCmd() {
	viper.GetViper().Set(util.CliTransportParamK8sArgs, &genericclioptions.ConfigFlags{})
	s.Run("should get config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		doguConfigServiceFactoryMock.EXPECT().Delete(configKey).Return(nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"delete", configKey})
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
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		configKey := "redmineKey"
		expectedError := errors.New("configService error")
		doguConfigServiceFactoryMock.EXPECT().Delete(configKey).Return(expectedError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		configCmd.SetArgs([]string{"delete", doguName, configKey})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config delete <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, fmt.Sprintf("cannot delete config key '%s' in delete dogu config command: configService error", configKey))
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)
		doguName := "redmine"
		viper.GetViper().Set(util.CliTransportArgConfigDoguDoguName, doguName)

		configKey := "redmineKey"
		expectedError := errors.New("create configService error")
		doguConfigServiceFactory = errorDoguConfigServiceFactory(expectedError)

		// when
		configCmd.SetArgs([]string{"delete", doguName, configKey})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config delete <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "cannot create config service in get dogu config command: create configService error")
	})

	s.Run("should fail with too few Arguments", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		configCmd := Cmd()
		configCmd.SetOut(outBuf)
		configCmd.SetErr(errBuf)

		// when
		configCmd.SetArgs([]string{"delete"})
		err := configCmd.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:\n  config delete <doguName> <configKey> [flags]", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts 2 arg(s), received 0")
	})
}

func noopDoguConfigServiceFactory(configServiceMock *mockDoguConfigService) func(doguName string, k8sNamespace string, restConfig *rest.Config) (doguConfigService, error) {
	return func(doguName, k8sNamespace string, restConfig *rest.Config) (doguConfigService, error) {
		return configServiceMock, nil
	}
}

func errorDoguConfigServiceFactory(expectedError error) func(doguName string, k8sNamespace string, restConfig *rest.Config) (doguConfigService, error) {
	return func(doguName, k8sNamespace string, restConfig *rest.Config) (doguConfigService, error) {
		return nil, expectedError
	}
}

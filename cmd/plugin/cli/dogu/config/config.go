package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu/config"
)

const (
	errMsgDoguConfigServiceCreate = "cannot create config service in get dogu config command: %w"
)

var doguConfigServiceFactory = func(doguName, k8sNamespace string, restConfig *rest.Config) (doguConfigService, error) {
	return config.New(doguName, k8sNamespace, restConfig)
}

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c", "cfg", "conf"},
	}

	cmd.AddCommand(
		listAllForDoguCmd(),
		getCmd(),
		editCmd(),
		deleteCmd(),
	)

	return cmd
}

func listAllForDoguCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := getTransportArgAsString(util.CliTransportArgConfigDoguDoguName)
			namespace := "test-namespace"
			k8sArgs := getTransportArg(util.CliTransportParamK8sArgs)
			restConfig, err := createRestConfig(k8sArgs)
			if err != nil {
				return err
			}

			configService, err := doguConfigServiceFactory(doguName, namespace, restConfig)
			if err != nil {
				return fmt.Errorf(errMsgDoguConfigServiceCreate, err)
			}

			configEntries, err := configService.GetAllForDogu()
			if err != nil {
				return fmt.Errorf("cannot list config in list dogu config command: %w", err)
			}

			for key, value := range configEntries {
				cmd.Printf("%s: %s\n", key, value)
			}
			return nil
		},
	}

	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <configKey>",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]

			doguName := getTransportArgAsString(util.CliTransportArgConfigDoguDoguName)
			namespace := "test-namespace"
			k8sArgs := getTransportArg(util.CliTransportParamK8sArgs)
			restConfig, err := createRestConfig(k8sArgs)
			if err != nil {
				return err
			}

			configService, err := doguConfigServiceFactory(doguName, namespace, restConfig)
			if err != nil {
				return fmt.Errorf(errMsgDoguConfigServiceCreate, err)
			}

			configValue, err := configService.GetValue(configKey)
			if err != nil {
				return fmt.Errorf("cannot get config key '%s' in get dogu config command: %w", configKey, err)
			}

			cmd.Printf(configValue)
			return nil
		},
	}

	return cmd
}

func editCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit <configKey> <configValue>",
		Aliases: []string{"e", "set"},
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]
			configValue := args[2]

			doguName := getTransportArgAsString(util.CliTransportArgConfigDoguDoguName)
			namespace := "test-namespace"
			k8sArgs := getTransportArg(util.CliTransportParamK8sArgs)
			restConfig, err := createRestConfig(k8sArgs)
			if err != nil {
				return err
			}

			configService, err := doguConfigServiceFactory(doguName, namespace, restConfig)
			if err != nil {
				return fmt.Errorf(errMsgDoguConfigServiceCreate, err)
			}

			err = configService.Edit(configKey, configValue)
			if err != nil {
				return fmt.Errorf("cannot set config key '%s' in edit dogu config command: %w", configKey, err)
			}

			return nil
		},
	}

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <configKey>",
		Aliases: []string{"d", "remove", "rm"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]

			doguName := getTransportArgAsString(util.CliTransportArgConfigDoguDoguName)
			namespace := "test-namespace"
			k8sArgs := getTransportArg(util.CliTransportParamK8sArgs)
			restConfig, err := createRestConfig(k8sArgs)
			if err != nil {
				return err
			}

			configService, err := doguConfigServiceFactory(doguName, namespace, restConfig)
			if err != nil {
				return fmt.Errorf(errMsgDoguConfigServiceCreate, err)
			}

			err = configService.Delete(configKey)
			if err != nil {
				return fmt.Errorf("cannot delete config key '%s' in delete dogu config command: %w", configKey, err)
			}

			return nil
		},
	}

	return cmd
}

func createRestConfig(k8sArgs interface{}) (*rest.Config, error) {
	restConfig, err := (k8sArgs).(*genericclioptions.ConfigFlags).ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("could not create rest config: %w", err)
	}
	return restConfig, nil
}

func getTransportArgAsString(paramName string) string {
	return viper.GetViper().GetString(paramName)
}

func getTransportArg(paramName string) interface{} {
	return viper.GetViper().Get(paramName)
}

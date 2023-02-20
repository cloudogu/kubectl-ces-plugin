package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu/config"
)

func Cmd(k8sArgs *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c", "cfg", "conf"},
	}

	cmd.AddCommand(
		listAllForDoguCmd(k8sArgs),
		getCmd(k8sArgs),
		editCmd(k8sArgs),
		deleteCmd(k8sArgs),
	)

	return cmd
}

var DoguConfigServiceFactory = func(viper *viper.Viper, k8sArgs *genericclioptions.ConfigFlags) (doguConfigService, error) {
	// TODO: add real namespace
	doguName := viper.GetString("doguName")
	restConfig, err := k8sArgs.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("could not create rest config: %w", err)
	}

	service, err := config.New(doguName, "test-namespace", restConfig)
	return service, err
}

func listAllForDoguCmd(k8sArgs *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configService, err := DoguConfigServiceFactory(viper.GetViper(), k8sArgs)
			if err != nil {
				return fmt.Errorf("cannot create config service in list dogu config command: %w", err)
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

func getCmd(k8sArgs *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <configKey>",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]

			configService, err := DoguConfigServiceFactory(viper.GetViper(), k8sArgs)
			if err != nil {
				return fmt.Errorf("cannot create config service in get dogu config command: %w", err)
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

func editCmd(k8sArgs *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit <configKey> <configValue>",
		Aliases: []string{"e", "set"},
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]
			configValue := args[2]

			configService, err := DoguConfigServiceFactory(viper.GetViper(), k8sArgs)
			if err != nil {
				return fmt.Errorf("cannot create config service in get dogu config command: %w", err)
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

func deleteCmd(k8sArgs *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <configKey>",
		Aliases: []string{"d", "remove", "rm"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]

			configService, err := DoguConfigServiceFactory(viper.GetViper(), k8sArgs)
			if err != nil {
				return fmt.Errorf("cannot create config service in get dogu config command: %w", err)
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

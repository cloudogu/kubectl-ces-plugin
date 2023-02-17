package config

import (
	"fmt"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",
	}

	cmd.AddCommand(
		listAllForDoguCmd(),
		getCmd(),
		editCmd(),
		deleteCmd(),
	)

	return cmd
}

type DoguConfigService interface {
	Edit(doguName string, registryKey string, registryValue string) error
	Delete(doguName string, registryKey string) error
	GetAllForDogu(doguName string) (map[string]string, error)
	GetValue(doguName string, registryKey string) (string, error)
}

var DoguConfigServiceFactory = func(viper *viper.Viper) (DoguConfigService, error) {
	//TODO: add real namespace and Rest-Config
	service, err := config.NewPortForwardedDoguConfigService("test namespace", nil)
	return service, err
}

func listAllForDoguCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list <doguName>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configService, err := DoguConfigServiceFactory(viper.GetViper())
			if err != nil {
				return fmt.Errorf("cannot create config service in list dogu config command: %w", err)
			}

			configEntries, err := configService.GetAllForDogu(viper.GetString("doguName"))
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
		Use:  "get <doguName> <configKey>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := viper.GetString("doguName")
			configKey := args[1]

			configService, err := DoguConfigServiceFactory(viper.GetViper())
			if err != nil {
				return fmt.Errorf("cannot create config service in get dogu config command: %w", err)
			}

			configValue, err := configService.GetValue(doguName, configKey)
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
		Use: "edit <dogu-name> <configKey> <configValue>",
		// TODO add completion, dogu name validation
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := viper.GetString("doguName")
			configKey := args[1]
			configValue := args[2]

			configService, err := DoguConfigServiceFactory(viper.GetViper())
			if err != nil {
				return fmt.Errorf("cannot create config service in get dogu config command: %w", err)
			}

			err = configService.Edit(doguName, configKey, configValue)
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
		Use: "delete <dogu-name> <registry-key>",
		// TODO add completion, dogu name validation
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	return cmd
}

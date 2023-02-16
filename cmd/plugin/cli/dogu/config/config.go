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

var DoguConfigServiceFactory = func(viper *viper.Viper) (*config.DoguConfigService, error) {
	//TODO: add real namespace and Rest-Config
	service, err := config.NewDoguConfigService("test namespace", nil)
	return service, err
}

func listAllForDoguCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
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
				cmd.Printf("%s, %s\n", key, value)
			}
			return nil
		},
	}

	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[1]
			doguName := viper.GetString("doguName")
			cmd.Printf("exec dogu config get command. dogu name: %s, config key: %s\n", doguName, configKey)
			cmd.Printf("viper config: %v\n", viper.AllSettings())
			return nil
		},
	}

	return cmd
}

func editCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "edit <dogu-name> <registry-key> <value>",
		// TODO add completion, dogu name validation
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
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

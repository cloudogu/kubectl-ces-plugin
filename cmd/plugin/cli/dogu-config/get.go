package dogu_config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <dogu> <configKey>",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := ""
			configKey := ""

			switch len(args) {
			case 1:
				doguName = args[0]
			case 2:
				doguName = args[0]
				configKey = args[1]
			}

			configService, err := doguConfigServiceFactory(doguName)
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

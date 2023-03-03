package dogu_config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func deleteCmd(factory serviceFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <dogu> <config-key>",
		Long:    `Delete a given configuration key.`,
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := ""
			configKey := ""

			if len(args) == 2 {
				doguName = args[0]
				configKey = args[1]
			}

			configService, err := factory.create(doguName)
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

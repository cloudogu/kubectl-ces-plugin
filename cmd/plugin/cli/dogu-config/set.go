package dogu_config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func setCmd(factory serviceFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set <dogu> <config-key> <config-value>",
		Long:    `Create or update a value for a given configuration key.`,
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := ""
			configKey := ""
			configValue := ""

			switch len(args) {
			case 3:
				doguName = args[0]
				configKey = args[1]
				configValue = args[2]
			}

			configService, err := factory.create(doguName)
			if err != nil {
				return fmt.Errorf(errMsgDoguConfigServiceCreate, err)
			}

			err = configService.Set(configKey, configValue)
			if err != nil {
				return fmt.Errorf("cannot set config key '%s': %w", configKey, err)
			}

			return nil
		},
	}

	return cmd
}

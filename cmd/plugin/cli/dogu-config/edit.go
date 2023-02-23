package dogu_config

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	flagKeyDeleteOnEmptyLong  = "delete-on-empty"
	flagKeyDeleteOnEmptyShort = "d"
)

var (
	flagValueDeleteOnEmpty bool
)

func editCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("edit <dogu> [<config-key>] [-%s | --%s]", flagKeyDeleteOnEmptyShort, flagKeyDeleteOnEmptyLong),
		Aliases: []string{"e"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().BoolVarP(&flagValueDeleteOnEmpty, flagKeyDeleteOnEmptyLong, flagKeyDeleteOnEmptyShort, false,
				"delete key if no value was provided during editing")
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

			err = configService.Edit(configKey, flagValueDeleteOnEmpty)
			if err != nil {
				return fmt.Errorf("cannot edit config keys: %w", err)
			}

			return nil
		},
	}

	return cmd
}

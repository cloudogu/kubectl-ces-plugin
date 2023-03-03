package dogu_config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCmd(factory serviceFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <dogu>",
		Long:    `Fetch a list of possible configuration keys.`,
		Aliases: []string{"l", "ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := ""
			if len(args) == 1 {
				doguName = args[0]
			}

			configService, err := factory.create(doguName)
			if err != nil {
				return fmt.Errorf(errMsgDoguConfigServiceCreate, err)
			}

			configEntries, err := configService.List()
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

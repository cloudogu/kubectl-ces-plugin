package dogu_config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <dogu>",
		Aliases: []string{"l", "ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			doguName := ""
			switch len(args) {
			case 1:
				doguName = args[0]
			}

			configService, err := doguConfigServiceFactory(doguName)
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

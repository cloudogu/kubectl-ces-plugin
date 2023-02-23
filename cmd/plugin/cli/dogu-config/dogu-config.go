package dogu_config

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dogu-config { edit | list | get | set | delete } <dogu> [<args>]...",
		Aliases: []string{"dc"},
	}

	cmd.AddCommand(
		editCmd(),
		listCmd(),
		getCmd(),
		setCmd(),
		deleteCmd(),
	)
	return cmd
}

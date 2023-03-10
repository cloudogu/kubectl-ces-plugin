package dogu_config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dogu-config { edit | list | get | set | delete } <dogu> [<args>]...",
		Aliases: []string{"dc"},
		Short:   "Manage a dogu's configuration",
		Long:    "Edit, list, get, set or delete values from a dogu's registry configuration.",
	}

	factory := &defaultServiceFactory{
		cliConfig: viper.GetViper(),
	}
	cmd.AddCommand(
		editCmd(factory),
		listCmd(factory),
		getCmd(factory),
		setCmd(factory),
		deleteCmd(factory),
	)
	return cmd
}

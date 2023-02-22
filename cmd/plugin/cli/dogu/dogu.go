package dogu

import (
	"github.com/spf13/cobra"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dogu <dogu-name> <command>",
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// doguArg := args[0]
			// viper.Set(util.CliTransportArgConfigDoguDoguName, doguArg)
		},
	}
	cmd.AddCommand(config.Cmd())
	return cmd
}

package dogu

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dogu <command>",
		Aliases: []string{"d"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			doguArg := args[0]
			viper.Set(util.CliTransportArgConfigDoguDoguName, doguArg)
		},
	}
	cmd.AddCommand(config.Cmd())
	return cmd
}

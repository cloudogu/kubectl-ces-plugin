package dogu

import (
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dogu <dogu-name> <command>",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			//TODO: Check if this works
			doguArg := cmd.PersistentFlags().Arg(1)
			viper.Set("doguName", doguArg)
			cmd.Printf("PersistentPreRun Args are: %v\n", args)
		},
	}
	cmd.AddCommand(config.Cmd())
	return cmd
}

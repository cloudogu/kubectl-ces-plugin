package dogu

import (
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dogu <command>",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			doguArg := args[0]
			viper.Set("doguName", doguArg)
			cmd.Printf("PersistentPreRun Args are: %v\n", args)
		},
	}
	cmd.AddCommand(config.Cmd())
	return cmd
}

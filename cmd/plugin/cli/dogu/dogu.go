package dogu

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
)

func Cmd(k8sArgs *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dogu <command>",
		Aliases: []string{"d"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			doguArg := args[0]
			viper.Set("doguName", doguArg)
		},
	}
	cmd.AddCommand(config.Cmd(k8sArgs))
	return cmd
}

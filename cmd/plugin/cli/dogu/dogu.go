package dogu

import (
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dogu [args]...",
	}

	cmd.AddCommand(config.Cmd())

	return cmd
}

package dogu_config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dogu-config { edit | list | get | set | delete } <dogu> [<args>]...",
		Aliases: []string{"dc"},
	}

	k8sArgs := viper.GetViper().Get(util.CliTransportParamK8sArgs)
	cfg := (k8sArgs).(*genericclioptions.ConfigFlags)
	namespace := ""
	if cfg.Namespace != nil {
		namespace = *cfg.Namespace
	}
	factory := &defaultServiceFactory{
		namespace:   namespace,
		configFlags: cfg,
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

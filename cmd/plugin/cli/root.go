package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu-config"
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
)

func RootCmd() *cobra.Command {
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	flags := pflag.NewFlagSet("kubectl", pflag.ExitOnError)
	pflag.CommandLine = flags
	KubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeResouceBuilderFlags := genericclioptions.NewResourceBuilderFlags()

	cmd := &cobra.Command{
		Use:           "kubectl ces",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cmd.SetErr(streams.ErrOut)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)
	flags.AddFlagSet(cmd.PersistentFlags())
	KubernetesConfigFlags.AddFlags(flags)
	kubeResouceBuilderFlags.AddFlags(flags)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.Set(util.CliTransportParamK8sArgs, KubernetesConfigFlags)

	cmd.AddCommand(dogu_config.Cmd())

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}

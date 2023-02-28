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
	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"
)

const (
	flagKeyLogLevel = logger.LogLevelKey
)

func RootCmd() *cobra.Command {
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	flags := pflag.NewFlagSet("kubectl", pflag.ExitOnError)

	// Set the default log level here. Alternatively supplied values overwrite the default during the flag parsing.
	flagValueLogLevel := logger.LogLevelWarn
	flagLogLevelRef := &flagValueLogLevel
	flags.Var(flagLogLevelRef, flagKeyLogLevel, "define log level")

	pflag.CommandLine = flags
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeResouceBuilderFlags := genericclioptions.NewResourceBuilderFlags()

	cmd := &cobra.Command{
		Use:           "kubectl ces",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger.ConfigureLogger()
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
	kubernetesConfigFlags.AddFlags(flags)
	kubeResouceBuilderFlags.AddFlags(flags)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.Set(util.CliTransportParamK8sArgs, kubernetesConfigFlags)
	viper.Set(util.CliTransportLogLevel, flagLogLevelRef)

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

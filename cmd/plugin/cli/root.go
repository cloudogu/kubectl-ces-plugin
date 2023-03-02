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
	flags := pflag.NewFlagSet("kubectl-ces", pflag.ExitOnError)

	// Set the default log level here. Alternatively supplied values overwrite the default during the flag parsing.
	flagValueLogLevel := logger.LogLevelWarn
	flagLogLevelRef := &flagValueLogLevel
	flags.Var(flagLogLevelRef, flagKeyLogLevel, "define the log level")

	pflag.CommandLine = flags
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeResourceBuilderFlags := genericclioptions.NewResourceBuilderFlags()

	cmd := &cobra.Command{
		Use:   "ces",
		Short: "Manage the Cloudogu EcoSystem",
		Long: `Provides various functions to make the management of the Cloudogu EcoSystem easier.
Among others, this includes editing dogu configurations.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := logger.ConfigureLogger()
			if err != nil {
				return err
			}

			err = viper.BindPFlags(cmd.Flags())
			cmd.SetErr(streams.ErrOut)
			if err != nil {
				return err
			}
			return nil
		},
	}

	// Cobra doesn't have a way to specify a two word command (ie. "kubectl krew"), so set a custom usage template
	// with kubectl in it. Cobra will use this template for the root and all child commands.
	cmd.SetUsageTemplate(strings.NewReplacer(
		"{{.UseLine}}", "kubectl {{.UseLine}}",
		"{{.CommandPath}}", "kubectl {{.CommandPath}}").Replace(cmd.UsageTemplate()))

	cobra.OnInitialize(initConfig)
	flags.AddFlagSet(cmd.PersistentFlags())
	kubernetesConfigFlags.AddFlags(flags)
	kubeResourceBuilderFlags.AddFlags(flags)

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

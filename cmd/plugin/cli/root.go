package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tj/go-spin"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubectl ces",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			//TODO: unhandled error
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger()
			log.Info("")

			s := spin.New()
			finishedCh := make(chan bool, 1)
			namespaceName := make(chan string, 1)
			go func() {
				lastNamespaceName := ""
				for {
					select {
					case <-finishedCh:
						fmt.Printf("\r")
						return
					case n := <-namespaceName:
						lastNamespaceName = n
					case <-time.After(time.Millisecond * 100):
						if lastNamespaceName == "" {
							fmt.Printf("\r  \033[36mSearching for namespaces\033[m %s", s.Next())
						} else {
							fmt.Printf("\r  \033[36mSearching for namespaces\033[m %s (%s)", s.Next(), lastNamespaceName)
						}
					}
				}
			}()
			defer func() {
				finishedCh <- true
			}()

			if err := plugin.RunPlugin(KubernetesConfigFlags, namespaceName); err != nil {
				return errors.Unwrap(err)
			}

			log.Info("")

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	KubernetesConfigFlags.AddFlags(cmd.Flags())

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cmd.AddCommand(dogu.Cmd())

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

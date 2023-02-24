package dogu_config

import (
	"fmt"

	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu-config"
)

const (
	errMsgDoguConfigServiceCreate = "cannot create config service: %w"
)

var doguConfigServiceFactory = func(doguName string) (doguConfigService, error) {
	k8sArgs := getTransportArg(util.CliTransportParamK8sArgs)
	restConfig, namespace, err := getKubeConfig(k8sArgs)
	if err != nil {
		return nil, err
	}

	return dogu_config.New(doguName, namespace, restConfig)
}

func getKubeConfig(k8sArgs interface{}) (*rest.Config, string, error) {
	cfg := (k8sArgs).(*genericclioptions.ConfigFlags)
	restConfig, err := cfg.ToRESTConfig()
	if err != nil {
		return nil, "", fmt.Errorf("could not create rest config: %w", err)
	}

	namespace := ""
	if cfg.Namespace != nil {
		namespace = *cfg.Namespace
	}

	return restConfig, namespace, nil
}

func getTransportArg(paramName string) interface{} {
	return viper.GetViper().Get(paramName)
}

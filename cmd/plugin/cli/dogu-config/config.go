package dogu_config

import (
	"fmt"
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/util"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu-config"
)

const (
	errMsgDoguConfigServiceCreate = "cannot create config service: %w"
)

type defaultServiceFactory struct {
	cliConfig configTransporter
}

func (sf *defaultServiceFactory) create(doguName string) (doguConfigService, error) {
	k8sArgs := sf.cliConfig.Get(util.CliTransportParamK8sArgs)
	cfg := (k8sArgs).(*genericclioptions.ConfigFlags)
	namespace := ""
	if cfg.Namespace != nil {
		namespace = *cfg.Namespace
	}

	restConfig, err := cfg.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("could not create rest config: %w", err)
	}

	return dogu_config.New(doguName, namespace, restConfig)
}

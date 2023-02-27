package dogu_config

import (
	"fmt"

	"github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu-config"
)

const (
	errMsgDoguConfigServiceCreate = "cannot create config service: %w"
)

type defaultServiceFactory struct {
	namespace   string
	configFlags restClientGetter
}

func (sf *defaultServiceFactory) create(doguName string) (doguConfigService, error) {
	restConfig, err := sf.configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("could not create rest config: %w", err)
	}

	return dogu_config.New(doguName, sf.namespace, restConfig)
}

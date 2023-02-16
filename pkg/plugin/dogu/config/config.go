package config

import (
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	"github.com/phayes/freeport"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
)

func NewDoguConfigService(namespace string, restConfig *rest.Config) (*DoguConfigService, error) {
	freePort, err := freeport.GetFreePort()

	endpoint := fmt.Sprintf("http://localhost:%d", freePort)
	reg, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, err
	}

	return &DoguConfigService{
		registry: reg,
		portForwarder: KubernetesPortForwarder{
			RestConfig: restConfig,
			Type:       ServiceType,
			NamespacedName: types.NamespacedName{
				Namespace: namespace,
				Name:      "etcd",
			},
			LocalPort:   freePort,
			ClusterPort: 4001,
		},
	}, nil
}

type DoguConfigService struct {
	registry      registry.Registry
	portForwarder PortForwarder
}

func (s DoguConfigService) Edit(doguName string, registryKey string, registryValue string) error {
	return s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		err = s.registry.DoguConfig(doguName).Set(registryKey, registryValue)
		if err != nil {
			return fmt.Errorf("error while editing key '%s' for dogu '%s': %w", registryKey, doguName, err)
		}

		return nil
	})
}

func (s DoguConfigService) Delete(doguName string, registryKey string) error {
	return s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		err = s.registry.DoguConfig(doguName).Delete(registryKey)
		if err != nil {
			return fmt.Errorf("error while deleting key '%s' for dogu '%s': %w", registryKey, doguName, err)
		}

		return nil
	})
}

func (s DoguConfigService) GetAllForDogu(doguName string) (map[string]string, error) {
	var configEntries map[string]string
	err := s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		configEntries, err = s.registry.DoguConfig(doguName).GetAll()
		if err != nil {
			return fmt.Errorf("error while reading all keys for dogu '%s': %w", doguName, err)
		}

		return nil
	})
	return configEntries, err
}

func (s DoguConfigService) GetValue(doguName string, registryKey string) (string, error) {
	var configValue string
	err := s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		configValue, err = s.registry.DoguConfig(doguName).Get(registryKey)
		if err != nil {
			return fmt.Errorf("error while reading key '%s' for dogu '%s': %w", registryKey, doguName, err)
		}

		return nil
	})
	return configValue, err
}

func (s DoguConfigService) checkInstallStatus(wantedDogu string) error {
	enabled, err := s.registry.DoguRegistry().IsEnabled(wantedDogu)
	if err != nil {
		return fmt.Errorf("cannot check if dogu '%s' is installed: %w", wantedDogu, err)
	}

	if !enabled {
		return fmt.Errorf("dogu '%s' is not installed", wantedDogu)
	}
	return nil
}

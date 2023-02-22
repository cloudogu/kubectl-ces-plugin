package config

import (
	"fmt"
	"strings"

	"github.com/phayes/freeport"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/portforward"
)

func New(doguName, namespace string, restConfig *rest.Config) (*doguConfigService, error) {
	forwarder, reg, err := createPortForwardAndRegistry(namespace, restConfig)
	if err != nil {
		return nil, err
	}

	delegator, err := newDelegator(doguName, forwarder, reg)
	if err != nil {
		return nil, err
	}

	return &doguConfigService{
		delegator: delegator,
	}, nil
}

func createPortForwardAndRegistry(namespace string, restConfig *rest.Config) (portForwarder, cesRegistry, error) {
	freePort, err := freeport.GetFreePort()
	if err != nil {
		return nil, nil, err
	}

	forward := portforward.New(restConfig, portforward.ServiceType, types.NamespacedName{Namespace: namespace, Name: "etcd"}, freePort, 4001)

	endpoint := fmt.Sprintf("http://localhost:%d", freePort)
	reg, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not create etcd registry: %w", err)
	}

	return forward, reg, nil
}

type doguConfigService struct {
	delegator delegator
}

func (s doguConfigService) Edit(registryKey string, registryValue string, deleteOnEmpty bool) error {
	return s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		if registryValue == "" {
			// edit interactive
			matchingFields := matchConfigurationFields(dogu, registryKey)
			if matchingFields == nil || len(matchingFields) == 0 {
				logger.NewLogger().Info("dogu '%s' has no matching configuration fields for key '%s'", dogu.GetSimpleName(), registryKey)
				return nil
			}

			return editor.EditConfiguration(matchingFields, deleteOnEmpty)
		} else {
			// set directly
			found, field := configurationFieldByKey(dogu, registryKey)
			if !found {
				logger.NewLogger().Info("dogu '%s' has no configuration field for key '%s'", dogu.GetSimpleName(), registryKey)
				return nil
			}

			return editor.SetFieldToValue(*field, registryValue, false)
		}
	})
}

func (s doguConfigService) Delete(registryKey string) error {
	return s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		found, field := configurationFieldByKey(dogu, registryKey)
		if !found {
			logger.NewLogger().Info("dogu '%s' has no configuration field for key '%s'", dogu.GetSimpleName(), registryKey)
			return nil
		}

		return editor.DeleteField(*field)
	})
}

func (s doguConfigService) GetAllForDogu() (map[string]string, error) {
	var entireDoguConfig map[string]string
	err := s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		for _, field := range dogu.Configuration {
			registryValue, err := editor.GetCurrentValue(field)
			if err != nil {
				return err
			}

			entireDoguConfig[field.Name] = registryValue
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return entireDoguConfig, nil
}

// GetValue returns the value of the provided registry key for the dogu if the key exists.
func (s doguConfigService) GetValue(registryKey string) (string, error) {
	var registryValue string
	err := s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		found, field := configurationFieldByKey(dogu, registryKey)
		if !found {
			logger.NewLogger().Info("dogu '%s' has no configuration field for key '%s'", dogu.GetSimpleName(), registryKey)
			return nil
		}

		var err error
		registryValue, err = editor.GetCurrentValue(*field)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return registryValue, nil
}

func configurationFieldByKey(dogu *core.Dogu, key string) (found bool, configField *core.ConfigurationField) {
	for _, field := range dogu.Configuration {
		if field.Name == key {
			return true, &field
		}
	}

	return false, nil
}

// matchConfigurationFields returns all configuration fields of a dogu who's names start with the provided key.
// An empty key returns all fields.
func matchConfigurationFields(dogu *core.Dogu, key string) []core.ConfigurationField {
	if key == "" {
		return dogu.Configuration
	}

	var matchingFields []core.ConfigurationField
	for _, field := range dogu.Configuration {
		if strings.HasPrefix(field.Name, key) {
			matchingFields = append(matchingFields, field)
		}
	}

	return matchingFields
}

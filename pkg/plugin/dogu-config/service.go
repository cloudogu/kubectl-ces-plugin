package dogu_config

import (
	"fmt"
	"strings"

	"github.com/phayes/freeport"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/portforward"
)

const noConfigFieldsErrorMsg = "dogu '%s' has no matching configuration fields for key '%s'"

func New(doguName, namespace string, restConfig *rest.Config) (*doguConfigService, error) {
	forwarder, reg, err := createPortForwardAndRegistry(namespace, restConfig)
	if err != nil {
		return nil, err
	}

	delegator := newDelegator(doguName, forwarder, reg)

	return &doguConfigService{
		delegator: delegator,
	}, nil
}

func createPortForwardAndRegistry(namespace string, restConfig *rest.Config) (portForwarder, cesRegistry, error) {
	freePort, err := freeport.GetFreePort()
	if err != nil {
		return nil, nil, fmt.Errorf("could not find free port for port-forward: %w", err)
	}

	forward := portforward.New(restConfig, types.NamespacedName{Namespace: namespace, Name: "etcd-0"}, freePort, 2379)

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

// Edit opens an interactive dialogue to edit the value of the given key.
func (s doguConfigService) Edit(registryKey string, deleteOnEmpty bool) error {
	return s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		matchingFields := matchConfigurationFields(dogu, registryKey)
		if len(matchingFields) == 0 {
			return fmt.Errorf(noConfigFieldsErrorMsg, dogu.GetSimpleName(), registryKey)
		}

		return editor.EditConfiguration(matchingFields, deleteOnEmpty)
	})
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

// Set replaces the current value of the key with the given value.
func (s doguConfigService) Set(registryKey string, registryValue string) error {
	return s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		found, field := configurationFieldByKey(dogu, registryKey)
		if !found {
			return fmt.Errorf(noConfigFieldsErrorMsg, dogu.GetSimpleName(), registryKey)
		}

		return editor.SetFieldToValue(*field, registryValue, false)
	})
}

// Delete deletes the given key and its value.
func (s doguConfigService) Delete(registryKey string) error {
	return s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		found, field := configurationFieldByKey(dogu, registryKey)
		if !found {
			return fmt.Errorf(noConfigFieldsErrorMsg, dogu.GetSimpleName(), registryKey)
		}

		return editor.DeleteField(*field)
	})
}

// List returns all keys and their values.
func (s doguConfigService) List() (map[string]string, error) {
	entireDoguConfig := map[string]string{}
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
			return fmt.Errorf(noConfigFieldsErrorMsg, dogu.GetSimpleName(), registryKey)
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

package config

import (
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
)

func New(doguName, namespace string, restConfig *rest.Config) (*doguConfigService, error) {
	delegator, err := newDelegator(doguName, namespace, restConfig)
	if err != nil {
		return nil, err
	}

	return &doguConfigService{
		delegator: delegator,
	}, nil
}

type doguConfigService struct {
	delegator *doguConfigurationDelegator
}

func (s doguConfigService) EditAllInteractive() error {
	return s.delegator.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
		return editor.EditConfiguration(dogu.Configuration)
	})
}

func (s doguConfigService) Edit(registryKey string, registryValue string) error {
	panic("todo")
}

func (s doguConfigService) Delete(registryKey string) error {
	panic("todo")
}

func (s doguConfigService) GetAllForDogu() (map[string]string, error) {
	panic("todo")
}

func (s doguConfigService) GetValue(registryKey string) (string, error) {
	panic("todo")
}

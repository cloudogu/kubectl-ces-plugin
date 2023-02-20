package config

import (
	"fmt"

	"github.com/phayes/freeport"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
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

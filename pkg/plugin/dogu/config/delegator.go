package config

import (
	"fmt"

	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"

	"github.com/phayes/freeport"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/doguConf"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/keys"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/portforward"
)

type doguConfigurationDelegator struct {
	doguName  string
	forwarder portForwarder
	doguReg   registry.DoguRegistry
	editor    doguConfigurationEditor
}

func newDelegator(doguName string, namespace string, restConfig *rest.Config) (*doguConfigurationDelegator, error) {
	freePort, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}

	forward := portforward.New(restConfig, portforward.ServiceType, types.NamespacedName{Namespace: namespace, Name: "etcd"}, freePort, 4001)

	endpoint := fmt.Sprintf("http://localhost:%d", freePort)
	etcd, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, fmt.Errorf("could not create etcd registry: %w", err)
	}

	var editor doguConfigurationEditor
	err = forward.ExecuteWithPortForward(func() error {

		keyManager, err := keys.NewKeyManager(etcd, doguName)
		if err != nil {
			return fmt.Errorf("could not create key manager for dogu '%s': %w", doguName, err)
		}

		publicKey, err := keyManager.GetPublicKey()
		if err != nil {
			return fmt.Errorf("could not get public key for dogu '%s': %w", doguName, err)
		}

		editor, err = doguConf.NewDoguConfigurationEditor(etcd.DoguConfig(doguName), publicKey)
		if err != nil {
			return fmt.Errorf("could not create configuration editor for dogu '%s': %w", doguName, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &doguConfigurationDelegator{
		doguName:  doguName,
		forwarder: forward,
		doguReg:   etcd.DoguRegistry(),
		editor:    editor,
	}, nil
}

// Delegate executes the given dogu configuration call against a previously configured registry.
func (dcd *doguConfigurationDelegator) Delegate(doguConfigCall func(dogu *core.Dogu, editor doguConfigurationEditor) error) error {
	err := dcd.forwarder.ExecuteWithPortForward(func() error {

		dogu, err := dcd.doguReg.Get(dcd.doguName)
		if err != nil {
			return fmt.Errorf("could not get dogu '%s' from etcd: %w", dcd.doguName, err)
		}

		if !doguConf.HasConfiguration(dogu) {
			logger.NewLogger().Info("dogu %s has no configuration fields", dogu.GetSimpleName())
			return nil
		}

		err = doguConfigCall(dogu, dcd.editor)
		if err != nil {
			return fmt.Errorf("error during registry interaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

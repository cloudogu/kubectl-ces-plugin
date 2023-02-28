package dogu_config

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/doguConf"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/keys"
)

type doguConfigurationDelegator struct {
	doguName    string
	forwarder   portForwarder
	doguReg     doguRegistry
	editFactory editorFactory
}

type defaultEditorFactory struct {
	doguName          string
	registry          cesRegistry
	keyManagerFactory keyManagerFactory
}

type defaultKeyManagerFactory struct{}

func (kmf *defaultKeyManagerFactory) create(registry cesRegistry, doguName string) (keyManager, error) {
	manager, err := keys.NewKeyManager(registry, doguName)
	if err != nil {
		return nil, fmt.Errorf("could not create key manager for dogu '%s': %w", doguName, err)
	}

	return manager, nil
}

func (ef *defaultEditorFactory) create() (doguConfigurationEditor, error) {
	keyMgr, err := ef.keyManagerFactory.create(ef.registry, ef.doguName)
	if err != nil {
		return nil, fmt.Errorf("could not create key manager for dogu '%s': %w", ef.doguName, err)
	}

	publicKey, err := keyMgr.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("could not get public key for dogu '%s': %w", ef.doguName, err)
	}

	doguConfig := ef.registry.DoguConfig(ef.doguName)
	return doguConf.NewDoguConfigurationEditor(doguConfig, publicKey)
}

func newDelegator(doguName string, forwarder portForwarder, reg cesRegistry) delegator {
	keyMgrFactory := &defaultKeyManagerFactory{}
	edFactory := &defaultEditorFactory{
		doguName:          doguName,
		registry:          reg,
		keyManagerFactory: keyMgrFactory,
	}

	return &doguConfigurationDelegator{
		doguName:    doguName,
		forwarder:   forwarder,
		doguReg:     reg.DoguRegistry(),
		editFactory: edFactory,
	}
}

// Delegate executes the given dogu configuration call against a previously configured registry.
func (dcd *doguConfigurationDelegator) Delegate(doguConfigCall func(dogu *core.Dogu, editor doguConfigurationEditor) error) error {
	err := dcd.forwarder.ExecuteWithPortForward(func() error {

		dogu, err := dcd.doguReg.Get(dcd.doguName)
		if err != nil {
			return fmt.Errorf("could not get dogu '%s' from etcd: %w", dcd.doguName, err)
		}

		if !doguConf.HasConfiguration(dogu) {
			fmt.Printf("dogu '%s' has no configuration fields\n", dogu.GetSimpleName())
			return nil
		}

		editor, err := dcd.editFactory.create()
		if err != nil {
			return fmt.Errorf("could not create configuration editor for dogu '%s': %w", dogu.GetSimpleName(), err)
		}

		err = doguConfigCall(dogu, editor)
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

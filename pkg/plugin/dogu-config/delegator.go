package dogu_config

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/doguConf"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/keys"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"
)

type doguConfigurationDelegator struct {
	doguName  string
	forwarder portForwarder
	doguReg   doguRegistry
	edFactory editorFactory
}

type editorFactory interface {
	create() (doguConfigurationEditor, error)
}

type defaultEditorFactory struct {
	doguName string
	registry cesRegistry
}

func (ef *defaultEditorFactory) create() (doguConfigurationEditor, error) {

	keyManager, err := keys.NewKeyManager(ef.registry, ef.doguName)
	if err != nil {
		return nil, fmt.Errorf("could not create key manager for dogu '%s': %w", ef.doguName, err)
	}

	publicKey, err := keyManager.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("could not get public key for dogu '%s': %w", ef.doguName, err)
	}

	editor, err := doguConf.NewDoguConfigurationEditor(ef.registry.DoguConfig(ef.doguName), publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not create configuration editor for dogu '%s': %w", ef.doguName, err)
	}

	return editor, nil
}

func newDelegator(doguName string, forwarder portForwarder, reg cesRegistry) delegator {
	edFactory := &defaultEditorFactory{
		doguName: doguName,
		registry: reg,
	}

	return &doguConfigurationDelegator{
		doguName:  doguName,
		forwarder: forwarder,
		doguReg:   reg.DoguRegistry(),
		edFactory: edFactory,
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
			logger.NewLogger().Info("dogu '%s' has no configuration fields", dogu.GetSimpleName())
			return nil
		}

		editor, err := dcd.edFactory.create()
		if err != nil {
			return err
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

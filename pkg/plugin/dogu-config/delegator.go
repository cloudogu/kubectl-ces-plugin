package dogu_config

import (
	"fmt"

	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/doguConf"
	"github.com/cloudogu/kubectl-ces-plugin/pkg/keys"
)

type doguConfigurationDelegator struct {
	doguName  string
	forwarder portForwarder
	doguReg   doguRegistry
	editor    doguConfigurationEditor
}

var newDelegator = func(doguName string, forwarder portForwarder, reg cesRegistry) (delegator, error) {
	var editor doguConfigurationEditor
	err := forwarder.ExecuteWithPortForward(func() error {

		keyManager, err := keys.NewKeyManager(reg, doguName)
		if err != nil {
			return fmt.Errorf("could not create key manager for dogu '%s': %w", doguName, err)
		}

		publicKey, err := keyManager.GetPublicKey()
		if err != nil {
			return fmt.Errorf("could not get public key for dogu '%s': %w", doguName, err)
		}

		editor, err = doguConf.NewDoguConfigurationEditor(reg.DoguConfig(doguName), publicKey)
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
		forwarder: forwarder,
		doguReg:   reg.DoguRegistry(),
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
			logger.NewLogger().Info("dogu '%s' has no configuration fields", dogu.GetSimpleName())
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

package keys

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
)

var log = core.GetLogger()

// NewKeyManager creates a new KeyManager which is able to find and create keys for the given dogu
func NewKeyManager(registry cesRegistry, doguName string) (*etcdKeyManager, error) {
	keyType, err := registry.GlobalConfig().Get("key_provider")
	if err != nil {
		return nil, fmt.Errorf("could not create etcd key manager: could not read key type: %w", err)
	}

	keyProvider, err := keys.NewKeyProvider(keyType)
	if err != nil {
		return nil, fmt.Errorf("could not create etcd key manager: could not create key provider: %w", err)
	}

	return &etcdKeyManager{
		doguConfig:  registry.DoguConfig(doguName),
		keyProvider: keyProvider,
	}, nil
}

type etcdKeyManager struct {
	doguConfig  configurationContext
	keyProvider keyProvider
}

// GetPublicKey returns the public key from the dogus etcd directory.
func (ekm *etcdKeyManager) GetPublicKey() (*keys.PublicKey, error) {
	exists, err := ekm.ExistsPublicKey()
	if err != nil {
		return nil, err
	}

	if exists {
		return ekm.fetchPublicKey()
	}

	return nil, fmt.Errorf("could not find public key in configuration context: %w", err)
}

// ExistsPublicKey true if the public key exists in the configuration context
func (ekm *etcdKeyManager) ExistsPublicKey() (bool, error) {
	exists, err := ekm.doguConfig.Exists(registry.KeyDoguPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to check if public key exists: %w", err)
	}
	return exists, nil
}

func (ekm *etcdKeyManager) fetchPublicKey() (*keys.PublicKey, error) {
	log.Debugf("fetch public key from configuration context")
	publicKeyString, err := ekm.doguConfig.Get(registry.KeyDoguPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key from registry: %w", err)
	}

	publicKey, err := ekm.keyProvider.ReadPublicKeyFromString(publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key from string: %w", err)
	}

	return publicKey, nil
}

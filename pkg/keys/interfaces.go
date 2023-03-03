package keys

import (
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
)

type keyProvider interface {
	ReadPublicKeyFromString(publicKeyString string) (*keys.PublicKey, error)
}

// cesRegistry represents the main registry of a Cloudogu EcoSystem. The registry
// manage dogus, their configuration and their states.
type cesRegistry interface {
	// Registry is embedded because using our own interfaces easily lets us generate mocks.
	// Therefore, we shifted to this pragmatic approach, which is much easier to read and understand anyway.
	registry.Registry
}

// configurationContext is able to manage the configuration of a single context
type configurationContext interface {
	registry.ConfigurationContext
}

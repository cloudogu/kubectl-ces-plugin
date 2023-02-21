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
	// Previously this interface made use of our own configurationContext interface like so:
	// ```
	// type cesRegistry interface {
	//	// GlobalConfig returns a configurationContext for the global context
	//	GlobalConfig() configurationContext
	//	// HostConfig returns a configurationContext for the host context
	//	HostConfig(hostService string) configurationContext
	//	// DoguConfig returns a configurationContext for the given dogu
	//	DoguConfig(dogu string) configurationContext
	//	// State returns the state object for the given dogu
	//	State(dogu string) registry.State
	//	// DoguRegistry returns an object which is able to manage dogus
	//	DoguRegistry() registry.DoguRegistry
	//	// BlueprintRegistry to maintain a blueprint history
	//	BlueprintRegistry() configurationContext
	//	// RootConfig returns a WatchConfigurationContext for the root context
	//	RootConfig() registry.WatchConfigurationContext
	//	// GetNode returns all keys that are included in any path, packed as Node
	//	GetNode() (registry.Node, error)
	// }
	// ```
	// Sadly, Go is too stoopid to do recursive duck-typing.
	// Therefore, we shifted to this pragmatic approach, which is much easier to read and understand anyway.
	registry.Registry
}

// configurationContext is able to manage the configuration of a single context
type configurationContext interface {
	registry.ConfigurationContext
}

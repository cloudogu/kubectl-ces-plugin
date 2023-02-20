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
	// GlobalConfig returns a configurationContext for the global context
	GlobalConfig() configurationContext
	// HostConfig returns a configurationContext for the host context
	HostConfig(hostService string) configurationContext
	// DoguConfig returns a configurationContext for the given dogu
	DoguConfig(dogu string) configurationContext
	// State returns the state object for the given dogu
	State(dogu string) registry.State
	// DoguRegistry returns an object which is able to manage dogus
	DoguRegistry() registry.DoguRegistry
	// BlueprintRegistry to maintain a blueprint history
	BlueprintRegistry() configurationContext
	// RootConfig returns a WatchConfigurationContext for the root context
	RootConfig() configurationContext
	// GetNode returns all keys that are included in any path, packed as Node
	GetNode() (registry.Node, error)
}

// configurationContext is able to manage the configuration of a single context
type configurationContext interface {
	// Set sets a configuration value in current context
	Set(key, value string) error
	// SetWithLifetime sets a configuration value in current context with the given lifetime
	SetWithLifetime(key, value string, timeToLiveInSeconds int) error
	// Refresh resets the time to live of a key
	Refresh(key string, timeToLiveInSeconds int) error
	// Get returns a configuration value from the current context
	Get(key string) (string, error)
	// GetAll returns a map of key value pairs
	GetAll() (map[string]string, error)
	// Delete removes a configuration key and value from the current context
	Delete(key string) error
	// DeleteRecursive removes a configuration key or directory from the current context
	DeleteRecursive(key string) error
	// Exists returns true if configuration key exists in the current context
	Exists(key string) (bool, error)
	// RemoveAll remove all configuration keys
	RemoveAll() error
	// GetOrFalse return false and empty string when the configuration value does not exist.
	// Otherwise, return true and the configuration value, even when the configuration value is an empty string.
	GetOrFalse(key string) (bool, string, error)
}

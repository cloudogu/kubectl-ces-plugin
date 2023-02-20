package config

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
)

// portForwarder provides functionality to create a port-forward.
type portForwarder interface {
	// ExecuteWithPortForward wraps the given function into a port-forward.
	ExecuteWithPortForward(fn func() error) error
}

// doguConfigurationEditor is able to edit registry configuration values of a dogu.
type doguConfigurationEditor interface {
	// EditConfiguration prints registry keys to writer and read values from reader.
	EditConfiguration(fields []core.ConfigurationField) error
	// GetCurrentValue returns a value for a given ConfigurationField if it exists, otherwise it returns an error.
	GetCurrentValue(field core.ConfigurationField) (string, error)
	// SetFieldToValue set the Field as value into the editor.
	SetFieldToValue(field core.ConfigurationField, value string) error
}

type keyManager interface {
	// GetPublicKey returns a dogu's public key
	GetPublicKey() (*keys.PublicKey, error)
	// ExistsPublicKey returns true if a dogu's public key exist.
	ExistsPublicKey() (bool, error)
}

// doguRegistry manages dogus on a ecosystem
type doguRegistry interface {
	// Enable enables the given dogu
	Enable(dogu *core.Dogu) error
	// Register registeres the dogu on the registry
	Register(dogu *core.Dogu) error
	// Unregister unregisters the dogu on the registry
	Unregister(name string) error
	// Get returns the dogu which the given name
	Get(name string) (*core.Dogu, error)
	// GetAll returns all installed dogus
	GetAll() ([]*core.Dogu, error)
	// IsEnabled returns true if the dogu is installed
	IsEnabled(name string) (bool, error)
}

// cesRegistry represents the main registry of a Cloudogu EcoSystem. The registry
// manage dogus, their configuration and their states.
type cesRegistry interface {
	// GlobalConfig returns a ConfigurationContext for the global context
	GlobalConfig() registry.ConfigurationContext
	// HostConfig returns a ConfigurationContext for the host context
	HostConfig(hostService string) registry.ConfigurationContext
	// DoguConfig returns a ConfigurationContext for the given dogu
	DoguConfig(dogu string) registry.ConfigurationContext
	// State returns the state object for the given dogu
	State(dogu string) registry.State
	// DoguRegistry returns an object which is able to manage dogus
	DoguRegistry() registry.DoguRegistry
	// BlueprintRegistry to maintain a blueprint history
	BlueprintRegistry() registry.ConfigurationContext
	// RootConfig returns a WatchConfigurationContext for the root context
	RootConfig() registry.WatchConfigurationContext
	// GetNode returns all keys that are included in any path, packed as Node
	GetNode() (registry.Node, error)
}

package dogu_config

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
)

type editorFactory interface {
	create() (doguConfigurationEditor, error)
}

type keyManagerFactory interface {
	create(registry cesRegistry, doguName string) (keyManager, error)
}

type keyManager interface {
	// GetPublicKey returns the public key from the dogus etcd directory.
	GetPublicKey() (*keys.PublicKey, error)
}

// portForwarder provides functionality to create a port-forward.
type portForwarder interface {
	// ExecuteWithPortForward provides the given function `fn()` with a given port-forward. Unless errors exists during
	// the port-forward instantiation, `fn()` guarantees a ready-to-use connection.
	ExecuteWithPortForward(fn func() error) error
}

// doguConfigurationEditor is able to edit registry configuration values of a dogu.
type doguConfigurationEditor interface {
	// EditConfiguration prints registry keys to writer and read values from reader.
	EditConfiguration(fields []core.ConfigurationField, deleteOnEmpty bool) error
	// GetCurrentValue returns a value for a given ConfigurationField if it exists, otherwise it returns an error.
	GetCurrentValue(field core.ConfigurationField) (string, error)
	// SetFieldToValue set the Field as value into the editor.
	SetFieldToValue(field core.ConfigurationField, value string, deleteOnEmpty bool) error
	// DeleteField deletes the field.
	DeleteField(field core.ConfigurationField) error
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
	registry.Registry
}

type delegator interface {
	Delegate(doguConfigCall func(dogu *core.Dogu, editor doguConfigurationEditor) error) error
}

package config

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/keys"
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

package config

// doguConfigService provides functionality to view and edit dogu configurations in etcd.
type doguConfigService interface {
	// Edit replaces the current value of the key with the given value.
	Edit(registryKey string, registryValue string) error
	// Delete deletes the given key and its value.
	Delete(registryKey string) error
	// GetAllForDogu returns all keys and their values.
	GetAllForDogu() (map[string]string, error)
	// GetValue returns the value of the given key.
	GetValue(registryKey string) (string, error)
}

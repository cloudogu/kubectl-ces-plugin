package dogu_config

type serviceFactory interface {
	create(doguName string) (doguConfigService, error)
}

// doguConfigService provides functionality to view and edit dogu configurations in etcd.
type doguConfigService interface {
	// Edit opens an interactive dialogue to edit the value of the given key.
	Edit(registryKey string, deleteOnEmpty bool) error
	// Set replaces the current value of the key with the given value.
	Set(registryKey string, registryValue string) error
	// Delete deletes the given key and its value.
	Delete(registryKey string) error
	// List returns all keys and their values.
	List() (map[string]string, error)
	// GetValue returns the value of the given key.
	GetValue(registryKey string) (string, error)
}

type configTransporter interface {
	// Get returns a previously Set configuration value.
	Get(key string) interface{}
	// Set persists the given value which can later be addressed with the key.
	Set(key string, value interface{})
}

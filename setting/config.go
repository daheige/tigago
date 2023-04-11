package setting

// Config interface associated with reading/saving configuration files.
type Config interface {
	// Load config load
	Load() error

	// IsSet is set value
	IsSet(key string) bool

	// ReadSection read val by key,val must be a pointer
	ReadSection(key string, val interface{}) error

	// Store save config to file
	Store(path string) error
}

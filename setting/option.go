package setting

// Option config option
type Options struct {
	configFile string
	watchFile  bool
}

// Option for ConfigOption
type Option func(*Options)

// WithConfigFile set config filename
func WithConfigFile(configFile string) Option {
	return func(c *Options) {
		c.configFile = configFile
	}
}

// WithWatchFile watch file change
func WithWatchFile() Option {
	return func(c *Options) {
		c.watchFile = true
	}
}

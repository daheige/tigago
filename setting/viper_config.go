package setting

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// New create a config interface.
func New(opts ...Option) Config {
	c := &viperConfig{
		vp:       viper.New(),
		sections: make(map[string]interface{}, 20),
	}

	conf := &Options{
		configFile: "./app.yaml",
	}

	for _, o := range opts {
		o(conf)
	}

	c.configFile = conf.configFile
	c.watchFile = conf.watchFile

	return c
}

type viperConfig struct {
	vp         *viper.Viper
	configFile string
	watchFile  bool
	sections   map[string]interface{}
}

// Load load config
func (c *viperConfig) Load() error {
	if c.configFile == "" {
		return errors.New("config file is empty")
	}

	configDir, err := filepath.Abs(filepath.Dir(c.configFile))
	if err != nil {
		return err
	}

	// file ext
	ext := strings.TrimPrefix(filepath.Ext(c.configFile), ".")
	if ext == "" {
		ext = "yaml" // default yaml file
	}

	file := filepath.Base(c.configFile)
	c.vp.SetConfigName(file)
	c.vp.AddConfigPath(configDir)
	c.vp.SetConfigType(ext)
	err = c.vp.ReadInConfig()
	if err != nil {
		return err
	}

	if c.watchFile {
		c.watch()
	}

	return nil
}

// IsSet is set value
func (c *viperConfig) IsSet(key string) bool {
	return c.vp.IsSet(key)
}

// ReadSection read val by key,val must be a pointer
func (c *viperConfig) ReadSection(key string, val interface{}) error {
	err := c.vp.UnmarshalKey(key, val)
	if _, ok := c.sections[key]; !ok {
		c.sections[key] = val
	}

	return err
}

// Store save config to file
func (c *viperConfig) Store(path string) error {
	return c.vp.WriteConfigAs(path)
}

func (c *viperConfig) reload() error {
	for k, v := range c.sections {
		err := c.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// when the configuration file is changed, reload all the configured keys/values
func (c *viperConfig) watch() {
	go func() {
		c.vp.WatchConfig()
		c.vp.OnConfigChange(func(in fsnotify.Event) {
			log.Println("config file has changed...")
			err := c.reload()
			if err != nil {
				log.Printf("read all config section err:%s\n", err.Error())
			}
		})
	}()
}

package rtconfig

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config is viper wrapper.
type Config struct {
	viper *viper.Viper

	mu       sync.Mutex
	settings map[string]interface{}
	watchers map[string][]WatcherCallback
}

// New returns new Config instance.
func New() *Config {
	config := &Config{
		viper:    viper.GetViper(),
		settings: make(map[string]interface{}),
		watchers: make(map[string][]WatcherCallback),
	}

	config.viper.AutomaticEnv()

	return config
}

// SetDefaultConfigFile sets default config file.
func (c *Config) SetDefaultConfigFile(path, name, ext string) error {
	if err := c.SetConfigFile(path, name, ext); err != nil {
		return err
	}

	for key, val := range c.viper.AllSettings() {
		c.viper.SetDefault(key, val)
	}

	return nil
}

// SetConfigFile sets config file.
func (c *Config) SetConfigFile(path, name, ext string) error {
	c.viper.SetConfigType(ext)
	c.viper.AddConfigPath(path)
	c.viper.SetConfigName(name)

	if err := c.viper.MergeInConfig(); err != nil {
		return fmt.Errorf("can't merge base config: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.settings = c.viper.AllSettings()

	return nil
}

// WatchConfig starts config watching.
func (c *Config) WatchConfig() {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(c.onConfigChange)
}

// Get returns a config Value associated with the key in viper config.
func (c *Config) Get(key string) (Value, error) {
	if val := c.viper.Get(key); val != nil {
		return value{val}, nil
	}

	return nil, ErrNilVariable
}

// WatchVariable allows to set a callback func on a specific variable change in viper config.
func (c *Config) WatchVariable(key string, cb WatcherCallback) error {
	if key == "" {
		return ErrEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key = strings.ToLower(key)

	c.watchers[key] = append(c.watchers[key], cb)

	return nil
}

func (c *Config) onConfigChange(_ fsnotify.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, newValue := range c.viper.AllSettings() {
		watchers, ok := c.watchers[key]
		if !ok {
			continue
		}

		oldValue, ok := c.settings[key]
		if !ok || reflect.DeepEqual(oldValue, newValue) {
			continue
		}

		for _, watcher := range watchers {
			watcher(value{oldValue}, value{newValue})
		}

		c.settings[key] = newValue
	}
}

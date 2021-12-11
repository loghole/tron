// Package config contains viper config initialisation method.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

// ViperConfig implements Config interface for values.yaml.
type ViperConfig struct {
	viper *viper.Viper

	mu       sync.Mutex
	settings map[string]interface{}
	watchers map[string][]WatcherCallback
}

// NewViperConfig returns new viper config instance.
func NewViperConfig(opts *app.Options) (*ViperConfig, error) {
	config := &ViperConfig{
		viper:    viper.GetViper(),
		settings: make(map[string]interface{}),
		watchers: make(map[string][]WatcherCallback),
	}

	if opts.ConfigMap != nil {
		if err := config.viper.MergeConfigMap(opts.ConfigMap); err != nil {
			return nil, fmt.Errorf("can't merge config map: %w", err)
		}

		return config, nil
	}

	config.viper.AutomaticEnv()

	config.viper.SetConfigType(app.ValuesExt)
	config.viper.AddConfigPath(filepath.Join(app.DeploymentsDir, app.ValuesDir))
	config.viper.SetConfigName(app.ValuesBaseName)

	// Init default config values.
	if err := config.viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("can't merge base config: %w", err)
	}

	for key, val := range config.viper.AllSettings() {
		config.viper.SetDefault(key, val)
	}

	config.viper.SetConfigName(app.ParseNamespace(os.Getenv(app.NamespaceEnv)).ValuesName())

	// Init config values for current namespace.
	if err := config.viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("can't merge namespace config: %w", err)
	}

	config.settings = config.viper.AllSettings()

	if opts.RealtimeConfig {
		config.viper.WatchConfig()
		config.viper.OnConfigChange(config.onConfigChange)
	}

	return config, nil
}

// Get returns a config Value associated with the key in viper config.
func (c *ViperConfig) Get(key string) (Value, error) {
	if val := c.viper.Get(key); val != nil {
		return value{val}, nil
	}

	return nil, ErrNilVariable
}

// WatchVariable allows to set a callback func on a specific variable change in viper config.
func (c *ViperConfig) WatchVariable(key string, cb WatcherCallback) error {
	if key == "" {
		return ErrEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key = strings.ToLower(key)

	c.watchers[key] = append(c.watchers[key], cb)

	return nil
}

func (c *ViperConfig) onConfigChange(_ fsnotify.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, newValue := range c.viper.AllSettings() {
		oldValue, ok := c.settings[key]
		if !ok || oldValue == newValue {
			continue
		}

		watchers, ok := c.watchers[key]
		if !ok {
			continue
		}

		for _, watcher := range watchers {
			watcher(value{oldValue}, value{newValue})
		}

		c.settings[key] = newValue
	}
}

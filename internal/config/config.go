// Package config contains viper config initialisation method.
package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

// Init initialise viper config.
func Init(opts *app.Options) error {
	viper.AutomaticEnv()

	if opts.ConfigMap != nil {
		return fromConfigMap(opts.ConfigMap)
	}

	return base(opts.RealtimeConfig)
}

func base(watch bool) error {
	viper.SetConfigType(app.ValuesExt)
	viper.AddConfigPath(filepath.Join(app.DeploymentsDir, app.ValuesDir))
	viper.SetConfigName(app.ValuesBaseName)

	// Init default config values.
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("can't merge base config: %w", err)
	}

	for key, val := range viper.AllSettings() {
		viper.SetDefault(key, val)
	}

	viper.SetConfigName(app.ParseNamespace(viper.GetString(app.NamespaceEnv)).ValuesName())

	// Init config values for current namespace.
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("can't merge namespace config: %w", err)
	}

	if watch {
		viper.WatchConfig()
	}

	return nil
}

func fromConfigMap(v map[string]interface{}) error {
	if err := viper.MergeConfigMap(v); err != nil {
		return fmt.Errorf("can't merge config map: %w", err)
	}

	return nil
}

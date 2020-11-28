// Package config contains viper config initialisation method.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lissteron/simplerr"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

// Init initialise viper config.
func Init(opts *app.Options) error {
	viper.AutomaticEnv()

	if opts.ConfigMap != nil {
		return fromConfigMap(opts.ConfigMap)
	}

	return base()
}

func base() error {
	viper.SetConfigType(app.ValuesExt)
	viper.SetConfigName(app.ValuesBaseName)
	viper.AddConfigPath(filepath.Join(app.DeploymentsDir, app.ValuesDir))

	if err := viper.ReadInConfig(); err != nil {
		return simplerr.Wrap(err, "read ")
	}

	namespace := app.ParseNamespace(viper.GetString(app.NamespaceEnv))

	replacer, err := os.Open(namespace.ValuesPath())
	if err != nil {
		return fmt.Errorf("open values file = '%s' failed: %w", namespace.ValuesPath(), err)
	}

	if err := viper.MergeConfig(replacer); err != nil {
		return fmt.Errorf("merge config failed: %w", err)
	}

	return nil
}

func fromConfigMap(v map[string]interface{}) error {
	if err := viper.MergeConfigMap(v); err != nil {
		return fmt.Errorf("can't merge config map: %w", err)
	}

	return nil
}

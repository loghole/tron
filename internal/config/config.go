package config

import (
	"os"
	"path/filepath"

	"github.com/lissteron/simplerr"
	"github.com/spf13/viper"

	"github.com/loghole/tron/internal/app"
)

func Init() error {
	viper.AutomaticEnv()
	viper.SetConfigType(app.ValuesExt)
	viper.SetConfigName(app.ValuesBaseName)
	viper.AddConfigPath(filepath.Join(app.DeploymentsDir, app.ValuesDir))

	if err := viper.ReadInConfig(); err != nil {
		return simplerr.Wrap(err, "read ")
	}

	namespace := app.ParseNamespace(viper.GetString(app.NamespaceEnv))

	replacer, err := os.Open(namespace.ValuesPath())
	if err != nil {
		return simplerr.Wrapf(err, "open values file = '%s' failed", namespace.ValuesPath())
	}

	if err := viper.MergeConfig(replacer); err != nil {
		return simplerr.Wrap(err, "merge config failed")
	}

	return nil
}

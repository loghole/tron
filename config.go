package tron

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/loghole/tron/internal/app"
	"github.com/loghole/tron/rtconfig"
)

func initConfig(info *app.Info, opts *app.Options) error {
	err := rtconfig.GetConfig().SetDefaultConfigFile(
		filepath.Join(app.DeploymentsDir, app.ValuesDir),
		app.ValuesBaseName,
		app.ValuesExt,
	)
	if err != nil && (!os.IsNotExist(err) || info.Namespace == app.NamespaceLocal) {
		return fmt.Errorf("init base values file: %w", err)
	}

	err = rtconfig.GetConfig().SetConfigFile(
		filepath.Join(app.DeploymentsDir, app.ValuesDir),
		info.Namespace.ValuesName(),
		app.ValuesExt,
	)
	if err != nil && (!os.IsNotExist(err) || info.Namespace == app.NamespaceLocal) {
		return fmt.Errorf("init namespace values file: %w", err)
	}

	if opts.RealtimeConfig {
		rtconfig.GetConfig().WatchConfig()
	}

	return nil
}

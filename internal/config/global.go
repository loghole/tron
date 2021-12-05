package config

import (
	"sync"

	"github.com/loghole/tron/internal/app"
)

// nolint:gochecknoglobals // global config instance.
var (
	_globalConfig     Config = StubConfig{}
	_globalConfigMu   sync.RWMutex
	_globalConfigOnce sync.Once
)

// Init initialise viper config.
func Init(opts *app.Options) error {
	_globalConfigMu.Lock()
	defer _globalConfigMu.Unlock()

	var err error

	_globalConfigOnce.Do(func() {
		_globalConfig, err = NewViperConfig(opts)
	})

	return err
}

// Get returns a config Value associated with the key.
func Get(key string) (Value, error) {
	_globalConfigMu.RLock()
	defer _globalConfigMu.RUnlock()

	return _globalConfig.Get(key) // nolint:wrapcheck // need clean error.
}

// WatchVariable allows to set a callback func on a specific variable change.
func WatchVariable(key string, cb WatcherCallback) error {
	_globalConfigMu.RLock()
	defer _globalConfigMu.RUnlock()

	return _globalConfig.WatchVariable(key, cb) // nolint:wrapcheck // need clean error.
}

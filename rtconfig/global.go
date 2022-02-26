package rtconfig

import (
	"sync"
	"time"
)

// nolint:gochecknoglobals // global config instance.
var (
	_globalConfig   = New()
	_globalConfigMu sync.RWMutex
)

// SetConfig sets global config instance.
func SetConfig(config *Config) {
	_globalConfigMu.Lock()
	defer _globalConfigMu.Unlock()

	_globalConfig = config
}

// GetConfig gets global config instance.
func GetConfig() *Config {
	_globalConfigMu.Lock()
	defer _globalConfigMu.Unlock()

	return _globalConfig
}

// GetValue returns a config Value associated with the key.
func GetValue(key string) (Value, error) {
	return GetConfig().Get(key)
}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	if val, err := GetConfig().Get(key); err == nil {
		return val.String()
	}

	return ""
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Bool()
	}

	return false
}

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Int()
	}

	return 0
}

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 {
	val, _ := GetConfig().Get(key)

	return val.Int32()
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Int64()
	}

	return 0
}

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Uint()
	}

	return 0
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Uint32()
	}

	return 0
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Uint64()
	}

	return 0
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Float64()
	}

	return 0
}

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Time()
	}

	return time.Time{}
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration {
	if val, err := GetConfig().Get(key); err == nil {
		return val.Duration()
	}

	return 0
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func GetIntSlice(key string) []int {
	if val, err := GetConfig().Get(key); err == nil {
		return val.IntSlice()
	}

	return nil
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string {
	if val, err := GetConfig().Get(key); err == nil {
		return val.StringSlice()
	}

	return nil
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} {
	if val, err := GetConfig().Get(key); err == nil {
		return val.StringMap()
	}

	return nil
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string {
	if val, err := GetConfig().Get(key); err == nil {
		return val.StringMapString()
	}

	return nil
}

// WatchVariable allows to set a callback func on a specific variable change.
func WatchVariable(key string, cb WatcherCallback) error {
	return GetConfig().WatchVariable(key, cb)
}

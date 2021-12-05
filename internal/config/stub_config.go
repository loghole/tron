package config

// StubConfig implements Config.
type StubConfig struct{}

// Get returns nilValue and ErrConfigNotInitialized error.
func (c StubConfig) Get(_ string) (Value, error) {
	return nilValue{}, ErrConfigNotInitialized
}

// WatchVariable returns ErrConfigNotInitialized error.
func (c StubConfig) WatchVariable(_ string, _ WatcherCallback) error {
	return ErrConfigNotInitialized
}

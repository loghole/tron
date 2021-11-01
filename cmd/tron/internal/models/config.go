package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const TronConfigFilepath = "tron" + sep + "config.json"

type Config struct {
	LastVersionCheck time.Time
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Read() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get user config dir: %w", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, TronConfigFilepath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}

func (c *Config) Write() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get user config dir: %w", err)
	}

	path := filepath.Join(dir, TronConfigFilepath)

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	if err := os.WriteFile(path, data, os.ModePerm); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

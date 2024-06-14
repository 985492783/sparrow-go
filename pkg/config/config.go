package config

import (
	"github.com/BurntSushi/toml"
)

type SparrowConfig struct {
	SwitcherConfig *SwitcherConfig
	ConsoleConfig  *ConsoleConfig
}

type SwitcherConfig struct {
	Enabled bool
	Addr    string
}

type ConsoleConfig struct {
	Enabled bool
	Addr    string
}

type LoggerConfig struct {
}

func newConfig() *SparrowConfig {
	return &SparrowConfig{
		SwitcherConfig: &SwitcherConfig{
			Enabled: true,
			Addr:    ":9854",
		},
		ConsoleConfig: &ConsoleConfig{
			Enabled: true,
			Addr:    ":9800",
		},
	}
}

func LoadConfig(path string) (*SparrowConfig, error) {
	config := newConfig()
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

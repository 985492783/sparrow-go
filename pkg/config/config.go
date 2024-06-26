package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
)

type SparrowConfig struct {
	SwitcherConfig *SwitcherConfig
	ConsoleConfig  *ConsoleConfig
	Auth           map[string]*Auth
	AuthEnabled    bool
}

type SwitcherConfig struct {
	Enabled bool
	Addr    string
}

type ConsoleConfig struct {
	Enabled bool
	Addr    string
}

type Auth struct {
	Password string
	Permits  string
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
		Auth:        make(map[string]*Auth),
		AuthEnabled: false,
	}
}

func (config *SparrowConfig) Authority(username, password, permit string) error {
	auths := config.Auth
	value, ok := auths[username]
	if !ok || value.Password != password {
		return errors.New(fmt.Sprintf("auth: %s unauthorized", username))
	}
	if strings.Contains(value.Permits, "-"+permit) {
		return errors.New(fmt.Sprintf("auth: %s permit error", username))
	}
	if strings.Contains(value.Permits, permit) || strings.Contains(value.Permits, "*:*") {
		return nil
	}
	split := strings.Split(permit, ":")
	if strings.Contains(value.Permits, split[0]+":*") {
		return nil
	}
	return errors.New(fmt.Sprintf("auth: %s permit error", username))
}

func LoadConfig(path string) (*SparrowConfig, error) {
	config := newConfig()
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

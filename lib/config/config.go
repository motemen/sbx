package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/shlex"
)

type Config struct {
	Default  *ProjectConfig           `json:"default,omitempty"`
	Projects map[string]ProjectConfig `json:"projects,omitempty"`
}

type ProjectConfig struct {
	// Session configuration for specific project (or default).
	Session *SessionConfig `json:"session,omitempty"`
}

type SessionConfig struct {
	// A command to obtain scrapbox.io session.
	Command string `json:"command,omitempty"`

	// A constant value of scrapbox.io session.
	Value string `json:"value,omitempty"`
}

func GetSession(projectName string) (string, error) {
	conf, err := Load()
	if err != nil {
		return "", err
	}

	return conf.GetSession(projectName)
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "sbx", "config.json")
	b, err := os.ReadFile(configPath)
	if err == os.ErrNotExist {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(b, &config)

	return &config, err
}

func (conf *Config) GetSession(projectName string) (string, error) {
	p := conf.GetProjectConfig(projectName)
	if p.Session == nil {
		return "", nil
	}

	if p.Session.Command != "" {
		parts, err := shlex.Split(p.Session.Command)
		if err != nil {
			return "", err
		}

		b, err := exec.Command(parts[0], parts[1:]...).Output()
		return string(b), err
	}

	if p.Session.Value != "" {
		return p.Session.Value, nil
	}

	return "", nil
}

func (conf *Config) GetProjectConfig(name string) (p ProjectConfig) {
	if conf == nil {
		return
	}

	p, ok := conf.Projects[name]
	if !ok && conf.Default != nil {
		p = *conf.Default
	}

	return
}

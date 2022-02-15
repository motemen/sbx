package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/shlex"
	"github.com/motemen/sbx/lib/sbapi"
)

type Config struct {
	Default  *ProjectConfig           `json:"default,omitempty"`
	Projects map[string]ProjectConfig `json:"projects,omitempty"`
}

type ProjectConfig struct {
	// Session configuration for specific project (or default).
	Session *SessionConfig `json:"session,omitempty"`

	Origin  string            `json:"origin,omitempty"`
	Host    string            `json:"host,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type SessionConfig struct {
	// A command to obtain scrapbox.io session.
	Command string `json:"command,omitempty"`

	// A constant value of scrapbox.io session.
	Value string `json:"value,omitempty"`
}

func GetOptions(projectName string) ([]sbapi.Option, error) {
	options := []sbapi.Option{}

	projectConf, err := GetProjectConfig(projectName)
	if err != nil {
		return nil, err
	}

	session, err := GetSession(projectName)
	if err != nil {
		return nil, err
	}

	options = append(options, sbapi.WithSessionID(session))

	if projectConf.Origin != "" {
		options = append(options, sbapi.WithOrigin(projectConf.Origin))
	}

	if projectConf.Host != "" {
		options = append(options, sbapi.WithHost(projectConf.Host))
	}

	if projectConf.Headers != nil {
		options = append(options, sbapi.WithHeaders(projectConf.Headers))
	}

	return options, nil
}

func GetSession(projectName string) (string, error) {
	conf, err := Load()
	if err != nil {
		return "", err
	}

	return conf.GetSession(projectName)
}

func GetProjectConfig(projectName string) (*ProjectConfig, error) {
	conf, err := Load()
	if err != nil {
		return nil, err
	}

	p := conf.GetProjectConfig(projectName)
	return &p, nil
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
	if err != nil {
		return nil, fmt.Errorf("%s: %w", configPath, err)
	}

	return &config, nil
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

func (p ProjectConfig) GetSession() (string, error) {
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

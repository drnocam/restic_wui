package config

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const ConfigFileName = "config.ini"

type Config struct {
	ResticCommand string `ini:"restic_command"`
}

func GetConfigPath() (string, error) {
	// For simplicity, we'll store config in the same directory as the executable or current working directory
	// In a real app, this might be in UserConfigDir
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ConfigFileName), nil
}

func Load() (*Config, error) {
	cfg := &Config{
		ResticCommand: "restic", // Default
	}

	path, err := GetConfigPath()
	if err != nil {
		return cfg, err
	}

	// If file doesn't exist, return default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	iniFile, err := ini.Load(path)
	if err != nil {
		return cfg, err
	}

	err = iniFile.MapTo(cfg)
	return cfg, err
}

func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	iniFile := ini.Empty()
	err = iniFile.ReflectFrom(cfg)
	if err != nil {
		return err
	}

	return iniFile.SaveTo(path)
}

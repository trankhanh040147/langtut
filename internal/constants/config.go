package constants

import (
	"os"
	"path/filepath"
)

const (
	ConfigDirName  = ".config"
	AppName        = "langtut"
	ConfigFileName = "config.yaml"
)

// GetConfigDir returns the config directory path
func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ConfigDirName, AppName)
}

// GetConfigPath returns the full config file path
func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), ConfigFileName)
}

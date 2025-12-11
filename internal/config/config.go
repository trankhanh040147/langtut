package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/trankhanh040147/langtut/internal/constants"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	APIKey         string            `yaml:"api_key"`
	TargetLanguage string            `yaml:"target_language"`
	Presets        map[string]string `yaml:"presets"`
}

// Load loads the configuration from file
func Load() (*Config, error) {
	configPath := constants.GetConfigPath()
	configDir := constants.GetConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	cfg := &Config{
		Presets: make(map[string]string),
	}

	// If config file doesn't exist, return defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure presets map is initialized
	if cfg.Presets == nil {
		cfg.Presets = make(map[string]string)
	}

	return cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config) error {
	configPath := constants.GetConfigPath()
	configDir := constants.GetConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// PromptAPIKey prompts the user for API key if not set
func PromptAPIKey(cfg *Config) error {
	if cfg.APIKey != "" {
		return nil
	}

	// Check if stdout is a TTY
	stat, err := os.Stdout.Stat()
	if err != nil || (stat.Mode()&os.ModeCharDevice) == 0 {
		return fmt.Errorf("api_key not set and not in interactive mode")
	}

	fmt.Fprintf(os.Stderr, "Gemini API key not found. Please enter your API key: ")
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	cfg.APIKey = strings.TrimSpace(apiKey)
	if cfg.APIKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Save config with API key
	if err := Save(cfg); err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}

	return nil
}

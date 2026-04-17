package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Config holds the global application configuration
type Config struct {
	ListenAddr  string `json:"listen_addr"`
	ListenPort  int    `json:"listen_port"`
	LogLevel    string `json:"log_level"`
	DataDir     string `json:"data_dir"`
	EnableHTTPS bool   `json:"enable_https"`
	CertFile    string `json:"cert_file"`
	KeyFile     string `json:"key_file"`
	AdminUser   string `json:"admin_user"`
	AdminPass   string `json:"admin_pass"`
}

var (
	globalConfig *Config
	configMu     sync.RWMutex
	configPath   string
)

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ListenAddr:  "0.0.0.0",
		ListenPort:  16601,
		LogLevel:    "info",
		DataDir:     "./data",
		EnableHTTPS: false,
		AdminUser:   "admin",
		AdminPass:   "admin666",
	}
}

// LoadConfig reads configuration from the given file path.
// If the file does not exist, a default config is created and saved.
func LoadConfig(path string) (*Config, error) {
	configPath = path

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, creating default config", path)
		cfg := DefaultConfig()
		if err := SaveConfig(cfg); err != nil {
			return nil, err
		}
		setGlobalConfig(cfg)
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	setGlobalConfig(cfg)
	log.Printf("Config loaded from %s", path)
	return cfg, nil
}

// SaveConfig writes the given configuration to disk
func SaveConfig(cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfig returns the current global configuration (thread-safe)
func GetConfig() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return globalConfig
}

// UpdateConfig updates the global config and persists it to disk
func UpdateConfig(cfg *Config) error {
	setGlobalConfig(cfg)
	return SaveConfig(cfg)
}

func setGlobalConfig(cfg *Config) {
	configMu.Lock()
	defer configMu.Unlock()
	globalConfig = cfg
}

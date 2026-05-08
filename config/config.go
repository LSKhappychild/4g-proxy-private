package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Config represents the proxy configuration
type Config struct {
	// Proxy settings
	Proxy ProxyConfig `yaml:"proxy"`

	// MME backend settings
	MME MMEConfig `yaml:"mme"`

	// HTTP API settings
	API APIConfig `yaml:"api"`

	// Logging settings
	Logging LoggingConfig `yaml:"logging"`

	// Delay settings (in milliseconds)
	Delay DelaySettingsConfig `yaml:"delay"`
}

// DelaySettingsConfig represents delay settings for each signal type (in milliseconds)
type DelaySettingsConfig struct {
	Attach           int64 `yaml:"attach"`
	Detach           int64 `yaml:"detach"`
	TAU              int64 `yaml:"tau"`
	ServiceRequest   int64 `yaml:"serviceRequest"`
	UEContextRelease int64 `yaml:"ueContextRelease"`
	PDNConnectivity  int64 `yaml:"pdnConnectivity"`
	Handover         int64 `yaml:"handover"`
	HandoverRequired int64 `yaml:"handoverRequired"` // Specific delay for HandoverRequired messages
	HandoverNotify   int64 `yaml:"handoverNotify"`   // Specific delay for HandoverNotify messages
	Reset            int64 `yaml:"reset"`
	Paging           int64 `yaml:"paging"`
	Default          int64 `yaml:"default"`
}

// ProxyConfig represents proxy settings
type ProxyConfig struct {
	// Listen address for eNB connections
	ListenAddress string `yaml:"listenAddress"`
	ListenPort    int    `yaml:"listenPort"`

	// SCTP settings
	SCTPInitMsgMaxInStreams  int `yaml:"sctpInitMsgMaxInStreams"`
	SCTPInitMsgMaxOutStreams int `yaml:"sctpInitMsgMaxOutStreams"`
}

// MMEConfig represents MME backend settings
type MMEConfig struct {
	// MME address
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// APIConfig represents HTTP API settings
type APIConfig struct {
	// Enable HTTP API
	Enabled bool `yaml:"enabled"`

	// Listen address
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// LoggingConfig represents logging settings
type LoggingConfig struct {
	// Log level
	Level string `yaml:"level"`

	// Log S1AP messages
	LogS1AP bool `yaml:"logS1AP"`

	// Log NAS messages
	LogNAS bool `yaml:"logNAS"`

	// Verbose mode
	Verbose bool `yaml:"verbose"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Proxy: ProxyConfig{
			ListenAddress:            "0.0.0.0",
			ListenPort:               36412,
			SCTPInitMsgMaxInStreams:  2,
			SCTPInitMsgMaxOutStreams: 2,
		},
		MME: MMEConfig{
			Address: "127.0.0.1",
			Port:    36412,
		},
		API: APIConfig{
			Enabled: true,
			Address: "0.0.0.0",
			Port:    8080,
		},
		Logging: LoggingConfig{
			Level:   "info",
			LogS1AP: true,
			LogNAS:  true,
			Verbose: false,
		},
		Delay: DelaySettingsConfig{
			Attach:           0,
			Detach:           0,
			TAU:              0,
			ServiceRequest:   0,
			UEContextRelease: 0,
			PDNConnectivity:  0,
			Handover:         0,
			HandoverRequired: 0,
			HandoverNotify:   0,
			Reset:            0,
			Paging:           0,
			Default:          0,
		},
	}
}

// LoadFromEnv loads configuration from environment variables
// Environment variables take precedence over config file values
func (c *Config) LoadFromEnv() {
	// Proxy settings
	if v := os.Getenv("PROXY_LISTEN_ADDRESS"); v != "" {
		c.Proxy.ListenAddress = v
	}
	if v := os.Getenv("PROXY_LISTEN_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Proxy.ListenPort = port
		}
	}

	// MME settings
	if v := os.Getenv("MME_ADDRESS"); v != "" {
		c.MME.Address = v
	}
	if v := os.Getenv("MME_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.MME.Port = port
		}
	}

	// API settings
	if v := os.Getenv("API_ENABLED"); v != "" {
		c.API.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("API_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.API.Port = port
		}
	}

	// Delay settings (in milliseconds)
	if v := os.Getenv("DELAY_ATTACH_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.Attach = ms
		}
	}
	if v := os.Getenv("DELAY_DETACH_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.Detach = ms
		}
	}
	if v := os.Getenv("DELAY_TAU_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.TAU = ms
		}
	}
	if v := os.Getenv("DELAY_SERVICE_REQUEST_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.ServiceRequest = ms
		}
	}
	if v := os.Getenv("DELAY_UE_CONTEXT_RELEASE_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.UEContextRelease = ms
		}
	}
	if v := os.Getenv("DELAY_PDN_CONNECTIVITY_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.PDNConnectivity = ms
		}
	}
	if v := os.Getenv("DELAY_HANDOVER_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.Handover = ms
		}
	}
	if v := os.Getenv("DELAY_HANDOVER_REQUIRED_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.HandoverRequired = ms
		}
	}
	if v := os.Getenv("DELAY_HANDOVER_NOTIFY_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.HandoverNotify = ms
		}
	}
	if v := os.Getenv("DELAY_RESET_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.Reset = ms
		}
	}
	if v := os.Getenv("DELAY_PAGING_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.Paging = ms
		}
	}
	if v := os.Getenv("DELAY_DEFAULT_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.Delay.Default = ms
		}
	}
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to a YAML file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ProxyEndpoint returns the proxy listen endpoint
func (c *Config) ProxyEndpoint() string {
	return fmt.Sprintf("%s:%d", c.Proxy.ListenAddress, c.Proxy.ListenPort)
}

// MMEEndpoint returns the MME backend endpoint
func (c *Config) MMEEndpoint() string {
	return fmt.Sprintf("%s:%d", c.MME.Address, c.MME.Port)
}

// APIEndpoint returns the HTTP API endpoint
func (c *Config) APIEndpoint() string {
	return fmt.Sprintf("%s:%d", c.API.Address, c.API.Port)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Proxy.ListenPort <= 0 || c.Proxy.ListenPort > 65535 {
		return fmt.Errorf("invalid proxy listen port: %d", c.Proxy.ListenPort)
	}

	if c.MME.Port <= 0 || c.MME.Port > 65535 {
		return fmt.Errorf("invalid MME port: %d", c.MME.Port)
	}

	if c.API.Enabled {
		if c.API.Port <= 0 || c.API.Port > 65535 {
			return fmt.Errorf("invalid API port: %d", c.API.Port)
		}
	}

	return nil
}

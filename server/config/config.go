package config

import (
	"encoding/json"
	"fmt"
)

// ParseConfig parses the raw JSON configuration.
func ParseConfig(raw []byte) (Config, error) {
	config := Config{
		TLS: &ConfigTLS{},
	}
	err := json.Unmarshal(raw, &config)
	if err != nil {
		return config, fmt.Errorf("unmarshal config: %v", err)
	}
	return config, nil
}

type Config struct {
	Server ConfigServer `json:"server"`
	TLS    *ConfigTLS   `json:"tls"`
}

type ConfigServer struct {
	HttpAddress  string `json:"http_address"`
	HttpsAddress string `json:"https_address"`
}

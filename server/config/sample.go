package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// CreateSample creates a sample configuration file.
func CreateSample(path string) error {
	sample := Config{
		Server: ConfigServer{
			HttpAddress:  ":7000",
			HttpsAddress: ":7001",
		},
		TLS: &ConfigTLS{
			DomainNameServer: []string{},
			IP:               []string{},
			Certificates:     []*ConfigTLSPath{},
		},
	}
	raw, err := json.MarshalIndent(&sample, "", "    ")
	if err != nil {
		return fmt.Errorf("could not marshal sample config: %v", err)
	}
	err = os.WriteFile(path, raw, 0600)
	if err != nil {
		return fmt.Errorf("could not write sample config file: %v", err)
	}
	return nil
}

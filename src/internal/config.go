package internal

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"gopkg.in/yaml.v2"
)

// LoadConfig loads the application configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	// Read the configuration file
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML content into the Config struct
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateTLSConfig creates a TLS configuration using the provided certificates
func CreateTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Load CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}, nil
}

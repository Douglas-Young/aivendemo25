package main

import (
	"log"
	"path/filepath"

	"github.com/Douglas-Young/aivendemo/src/internal"
)

func main() {
	// Get project root for proper file path finding
	projectRoot, err := filepath.Abs("../../")
	if err != nil {
		log.Fatalf("Failed to get project root path: %v", err)
	}

	configPath := filepath.Join(projectRoot, "config.yaml")

	config, err := internal.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config.Kafka.TLSCert = filepath.Join(projectRoot, config.Kafka.TLSCert)
	config.Kafka.TLSKey = filepath.Join(projectRoot, config.Kafka.TLSKey)
	config.Kafka.TLSCA = filepath.Join(projectRoot, config.Kafka.TLSCA)

	// Start producing Kafka messages in a goroutine and consume messages to send to opensearch
	internal.ProduceKafkaMessages(config)
	internal.ConsumeKafkaMessages(config)
}

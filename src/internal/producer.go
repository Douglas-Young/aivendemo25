package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/IBM/sarama"
)

var pages = []string{"/home", "/products", "/checkout", "/about", "/contact"}
var actions = []string{"click", "view", "scroll", "purchase"}

// generateClickEvent generates a simulated click event
func generateClickEvent() ClickEvent {
	return ClickEvent{
		Timestamp: time.Now().Format(time.RFC3339),         // Create timestamp in proper format
		UserID:    fmt.Sprintf("user-%d", rand.Intn(1000)), // Generate random unique id
		Page:      pages[rand.Intn(len(pages))],            // Generate random page
		Action:    actions[rand.Intn(len(actions))],        // Generate random action
	}
}

// ProduceKafkaMessages produces sample clickstream events to Kafka
func ProduceKafkaMessages(config *Config) {
	// TLS configuration for Kafka
	tlsConfig, err := CreateTLSConfig(config.Kafka.TLSCert, config.Kafka.TLSKey, config.Kafka.TLSCA)
	if err != nil {
		log.Fatalf("Failed to create TLS config: %v", err)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.TLS.Enable = true
	kafkaConfig.Net.TLS.Config = tlsConfig

	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Produce 10 sample messages and send to kafka topic in aiven
	for i := 0; i < 10; i++ {
		event := generateClickEvent()
		message, _ := json.Marshal(event)
		msg := &sarama.ProducerMessage{
			Topic: config.Kafka.Topic,
			Value: sarama.StringEncoder(message),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Printf("Message sent to partition %d at offset %d", partition, offset)
		}

		// Dont spam
		time.Sleep(1 * time.Second)
	}
}

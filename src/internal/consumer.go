package internal

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// ConsumeKafkaMessages consumes messages from Kafka and indexes them into OpenSearch
func ConsumeKafkaMessages(config *Config) {
	// TLS configuration for Kafka
	tlsConfig, err := CreateTLSConfig(config.Kafka.TLSCert, config.Kafka.TLSKey, config.Kafka.TLSCA)
	if err != nil {
		log.Fatalf("Failed to create TLS config: %v", err)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConfig.Net.TLS.Enable = true
	kafkaConfig.Net.TLS.Config = tlsConfig

	// Consumer group
	consumerGroup, err := sarama.NewConsumerGroup(config.Kafka.Brokers, config.Kafka.ConsumerGroup, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}
	defer consumerGroup.Close()

	handler := ConsumerHandler{config}
	for {
		if err := consumerGroup.Consume(context.Background(), []string{config.Kafka.Topic}, &handler); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}
}

// ConsumerHandler handles Kafka messages
type ConsumerHandler struct {
	config *Config
}

func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// OpenSearch client
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{h.config.OpenSearch.URL},
		Username:  h.config.OpenSearch.User,
		Password:  h.config.OpenSearch.Password,
	})
	if err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Process Kafka messages and send to opensearch api
	for msg := range claim.Messages() {
		var event ClickEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal Kafka message: %v", err)
			continue
		}

		req := opensearchapi.IndexRequest{
			Index: h.config.OpenSearch.Index,
			Body:  strings.NewReader(string(msg.Value)),
		}
		res, err := req.Do(context.Background(), client)
		if err != nil {
			log.Printf("Failed to index message into OpenSearch: %v", err)
		} else {
			defer res.Body.Close()
			log.Printf("Indexed message into OpenSearch: %+v", event)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

package internal

type Config struct {
	Kafka struct {
		Brokers       []string `yaml:"brokers"`
		Topic         string   `yaml:"topic"`
		ConsumerGroup string   `yaml:"consumer_group"`
		TLSCert       string   `yaml:"tls_cert"`
		TLSKey        string   `yaml:"tls_key"`
		TLSCA         string   `yaml:"tls_ca"`
	} `yaml:"kafka"`
	OpenSearch struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Index    string `yaml:"index"`
	} `yaml:"opensearch"`
}

// ClickEvent represents the simulated user click event from the assignment
type ClickEvent struct {
	Timestamp string `json:"timestamp"` // ISO 8601 format
	UserID    string `json:"user_id"`   // Unique user ID
	Page      string `json:"page"`      // Web page URL
	Action    string `json:"action"`    // Examples of some possible user actions (e.g., click, view, scroll)
}

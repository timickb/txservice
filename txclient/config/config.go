package config

type AppConfig struct {
	AppPort     int
	KafkaPort   int
	TxClientNum int // Defines kafka topic to write requests
	KafkaHost   string
}

func NewDefault() *AppConfig {
	return &AppConfig{
		AppPort:     8088,
		KafkaPort:   9092,
		KafkaHost:   "localhost",
		TxClientNum: 1,
	}
}

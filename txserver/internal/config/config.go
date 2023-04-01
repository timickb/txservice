package config

type postgres struct {
	Host     string
	Name     string
	User     string
	Password string
	SslMode  string
	Port     int
}

type kafka struct {
	Host  string
	Topic string
	Port  int
}

type ServerConfig struct {
	AppPort        int
	TxClientsCount int
	Database       *postgres
	Kafka          *kafka
}

func NewDefault() *ServerConfig {
	return &ServerConfig{
		AppPort:        8080,
		TxClientsCount: 1,
		Database: &postgres{
			Host:     "localhost",
			Name:     "postgres",
			User:     "postgres",
			Password: "",
			Port:     5432,
			SslMode:  "disable",
		},
		Kafka: &kafka{
			Host:  "localhost",
			Topic: "default_chan",
			Port:  9092,
		},
	}
}

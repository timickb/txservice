package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/timickb/txclient/config"
	"github.com/timickb/txclient/internal/delivery/http"
	"github.com/timickb/txclient/internal/infrastructure/messaging"
	"github.com/timickb/txclient/internal/service"
	"os"
	"strconv"
)

func main() {
	logger := logrus.New()

	logFmt := new(logrus.TextFormatter)
	logFmt.FullTimestamp = true
	logFmt.TimestampFormat = "2006-01-02 15:04:05"
	logger.SetFormatter(logFmt)

	if err := mainNoExit(logger); err != nil {
		logger.Fatal(err)
	}
}

func mainNoExit(logger *logrus.Logger) error {
	cfg := config.NewDefault()
	parseConfigFromEnvironment(cfg)

	kafkaAddr := fmt.Sprintf("%s:%d", cfg.KafkaHost, cfg.KafkaPort)
	// Each tx client has separate msg channel,
	// so it needs a kafka topic with name tx_{client_number}
	kafkaTopic := fmt.Sprintf("tx_%d", cfg.TxClientNum)

	logger.Info("Kafka URL: ", kafkaAddr)
	logger.Info("Kafka topic: ", kafkaTopic)

	prod := messaging.NewProducer(kafkaAddr, kafkaTopic, logger)
	sender := service.NewSender(prod, logger)

	ctx := context.Background()
	srv := http.New(ctx, sender, logger)

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = srv.Stop()
				logger.Fatal(ctx.Err())
			}
		}
	}()

	logger.Info("Starting HTTP server on port ", cfg.AppPort)
	if err := srv.Run(cfg.AppPort); err != nil {
		return err
	}

	return nil
}

func parseConfigFromEnvironment(cfg *config.AppConfig) {
	if os.Getenv("KAFKA_HOST") != "" {
		cfg.KafkaHost = os.Getenv("KAFKA_HOST")
	}
	if os.Getenv("KAFKA_PORT") != "" {
		cfg.KafkaPort, _ = strconv.Atoi(os.Getenv("KAFKA_PORT"))
	}
	if os.Getenv("APP_PORT") != "" {
		cfg.AppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	}
	if os.Getenv("TX_CLIENT_NUM") != "" {
		cfg.TxClientNum, _ = strconv.Atoi(os.Getenv("TX_CLIENT_NUM"))
	}
}

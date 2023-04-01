package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/timickb/txserver/internal/config"
	"github.com/timickb/txserver/internal/infrastructure/messaging"
	"github.com/timickb/txserver/internal/infrastructure/storage/postgres"
	"github.com/timickb/txserver/internal/interfaces"
	"github.com/timickb/txserver/internal/usecase"
	"github.com/timickb/txserver/internal/worker"
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

func mainNoExit(logger interfaces.Logger) error {
	cfg := config.NewDefault()
	parseConfigFromEnvironment(cfg)

	connStr := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=%s port=%d password=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.SslMode,
		cfg.Database.Port,
		cfg.Database.Password)

	logger.Info("Initializing postgres connection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	logger.Info("Starting requests queue reading worker")

	kafkaUrl := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	topics := txTopicsList(cfg.TxClientsCount)
	logger.Info("Topics to read: ", topics)

	repo := postgres.New(db)
	service := usecase.NewTransaction(repo, logger)
	cons := messaging.NewConsumer(kafkaUrl, "tx_handlers", topics, logger)
	w := worker.New(cons, service, logger)

	ctx := context.Background()
	errChan := make(chan error)

	// Worker is responsible for receiving messages from Kafka
	go func() {
		if err := w.Start(ctx, errChan); err != nil {
			logger.Fatal(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			_ = w.Stop()
			return ctx.Err()
		case err := <-errChan:
			logger.Error(err.Error())
		}
	}
}

func txTopicsList(count int) []string {
	var list []string
	for i := 1; i <= count; i++ {
		list = append(list, fmt.Sprintf("tx_%d", i))
	}
	return list
}

func parseConfigFromEnvironment(cfg *config.ServerConfig) {
	if os.Getenv("DB_HOST") != "" {
		cfg.Database.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_USER") != "" {
		cfg.Database.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_NAME") != "" {
		cfg.Database.Name = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_SSL_MODE") != "" {
		cfg.Database.SslMode = os.Getenv("DB_SSL_MODE")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		cfg.Database.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_PORT") != "" {
		cfg.Database.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	}
	if os.Getenv("APP_PORT") != "" {
		cfg.AppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	}
	if os.Getenv("KAFKA_HOST") != "" {
		cfg.Kafka.Host = os.Getenv("KAFKA_HOST")
	}
	if os.Getenv("TX_CLIENTS_COUNT") != "" {
		cfg.TxClientsCount, _ = strconv.Atoi(os.Getenv("TX_CLIENTS_COUNT"))
	}
	if os.Getenv("KAFKA_PORT") != "" {
		cfg.Kafka.Port, _ = strconv.Atoi(os.Getenv("KAFKA_PORT"))
	}
}

package kafka

import (
	"context"
	"log"
	"time"

	"github.com/nattatorn-dev/go-healthcheck/domain/entities"
	"github.com/nattatorn-dev/go-healthcheck/domain/repositories"
	"github.com/segmentio/kafka-go"
)

type KafkaHealth struct {
	address string
	config  entities.CheckerConfig
}

func NewKafkaHealth(address string, config entities.CheckerConfig) repositories.HealthChecker {
	return &KafkaHealth{address: address, config: config}
}

func (k *KafkaHealth) CheckHealth() entities.HealthCheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), k.config.Timeout)
	defer cancel()
	start := time.Now()
	conn, err := kafka.DialContext(ctx, "tcp", k.address)
	duration := time.Since(start)
	if err != nil {
		log.Printf("Kafka health check failed: %v", err)
		return entities.HealthCheckResult{Status: entities.HealthStatusDown, Error: err, Duration: duration.String()}
	}
	defer conn.Close()
	log.Printf("Kafka health check successful: %s", duration.String())
	return entities.HealthCheckResult{Status: entities.HealthStatusUp, Duration: duration.String()}
}

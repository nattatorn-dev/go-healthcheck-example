package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nattatorn-dev/go-healthcheck/domain/entities"
	"github.com/nattatorn-dev/go-healthcheck/domain/repositories"
)

type RedisHealth struct {
	client *redis.Client
	config entities.CheckerConfig
}

func NewRedisHealth(client *redis.Client, config entities.CheckerConfig) repositories.HealthChecker {
	return &RedisHealth{client: client, config: config}
}

func (r *RedisHealth) CheckHealth() entities.HealthCheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	start := time.Now()
	_, err := r.client.Ping(ctx).Result()
	duration := time.Since(start)
	if err != nil {
		log.Printf("Redis health check failed: %v", err)
		return entities.HealthCheckResult{Status: entities.HealthStatusDown, Error: err, Duration: duration.String()}
	}
	log.Printf("Redis health check successful: %s", duration.String())
	return entities.HealthCheckResult{Status: entities.HealthStatusUp, Duration: duration.String()}
}

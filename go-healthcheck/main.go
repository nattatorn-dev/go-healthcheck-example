package main

import (
	"log"
	"net/http"
	"time"

	redispkg "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/nattatorn-dev/go-healthcheck/domain/entities"
	"github.com/nattatorn-dev/go-healthcheck/domain/services"
	"github.com/nattatorn-dev/go-healthcheck/infra/db"
	"github.com/nattatorn-dev/go-healthcheck/infra/external"
	"github.com/nattatorn-dev/go-healthcheck/infra/kafka"
	"github.com/nattatorn-dev/go-healthcheck/infra/redis"
	"github.com/nattatorn-dev/go-healthcheck/interfaces/handlers"
)

func main() {
	// Per-checker configurations
	dbConfig := entities.CheckerConfig{
		Timeout:  1 * time.Second,
		Interval: 1 * time.Second,
	}
	redisConfig := entities.CheckerConfig{
		Timeout:  1 * time.Second,
		Interval: 3 * time.Second,
	}
	kafkaConfig := entities.CheckerConfig{
		Timeout:  1 * time.Second,
		Interval: 10 * time.Second,
	}
	externalAPIConfig := entities.CheckerConfig{
		Timeout:  3 * time.Second,
		Interval: 1 * time.Minute,
	}

	// Initialize dependencies
	redisClient := redispkg.NewClient(&redispkg.Options{Addr: "localhost:6379"})
	redisHealth := redis.NewRedisHealth(redisClient, redisConfig)
	kafkaHealth := kafka.NewKafkaHealth("localhost:9092", kafkaConfig)

	// DSN without database name
	dsn := "root:password@tcp(localhost:3306)/"
	dbHealth, err := db.NewSQLDBHealth("mysql", dsn, dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database health check: %v", err)
	}

	externalAPI := external.NewExternalAPIHealth("https://example.com", externalAPIConfig)

	// Initialize health store
	healthStore := entities.NewHealthStore()

	// Initialize service
	healthService := services.NewHealthService(healthStore)
	healthService.RegisterReadiness("Redis", redisHealth, redisConfig)
	healthService.RegisterReadiness("Kafka", kafkaHealth, kafkaConfig)
	healthService.RegisterReadiness("Database SQL", dbHealth, dbConfig)
	healthService.RegisterReadiness("External API", externalAPI, externalAPIConfig)

	healthService.RegisterLiveness("Redis", redisHealth, redisConfig)
	healthService.RegisterLiveness("Kafka", kafkaHealth, kafkaConfig)
	healthService.RegisterLiveness("Database SQL", dbHealth, dbConfig)
	healthService.RegisterLiveness("External API", externalAPI, externalAPIConfig)

	// Start background worker for health checks
	healthService.StartBackgroundCheck()

	// Initialize handlers
	readinessHandler := handlers.NewReadinessHandler(healthService)
	livenessHandler := handlers.NewLivenessHandler(healthService)

	// Start server
	http.Handle("/health/readiness", readinessHandler)
	http.Handle("/health/liveness", livenessHandler)
	log.Println("Starting server on :8080")
	http.ListenAndServe("127.0.0.1:8080", nil)
}

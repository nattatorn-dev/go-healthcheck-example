package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/nattatorn-dev/go-healthcheck/domain/entities"
	"github.com/nattatorn-dev/go-healthcheck/domain/repositories"
)

type SQLDBHealth struct {
	db     *sql.DB
	dsn    string
	config entities.CheckerConfig
}

func NewSQLDBHealth(driverName, dsn string, config entities.CheckerConfig) (repositories.HealthChecker, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &SQLDBHealth{db: db, dsn: dsn, config: config}, nil
}

func (d *SQLDBHealth) CheckHealth() entities.HealthCheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), d.config.Timeout)
	defer cancel()
	start := time.Now()
	err := d.db.PingContext(ctx)
	duration := time.Since(start)
	if err != nil {
		log.Printf("Database health check failed: %v", err)
		return entities.HealthCheckResult{Status: entities.HealthStatusDown, Error: err, Duration: duration.String()}
	}
	log.Printf("Database health check successful: %s", duration.String())
	return entities.HealthCheckResult{Status: entities.HealthStatusUp, Duration: duration.String()}
}

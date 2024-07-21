package repositories

import "github.com/nattatorn-dev/go-healthcheck/domain/entities"

type HealthChecker interface {
	CheckHealth() entities.HealthCheckResult
}

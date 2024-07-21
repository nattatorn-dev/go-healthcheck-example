package external

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nattatorn-dev/go-healthcheck/domain/entities"
	"github.com/nattatorn-dev/go-healthcheck/domain/repositories"
)

type ExternalAPIHealth struct {
	URL    string
	config entities.CheckerConfig
}

func NewExternalAPIHealth(url string, config entities.CheckerConfig) repositories.HealthChecker {
	return &ExternalAPIHealth{URL: url, config: config}
}

func (e *ExternalAPIHealth) CheckHealth() entities.HealthCheckResult {
	client := &http.Client{Timeout: e.config.Timeout}
	start := time.Now()
	resp, err := client.Get(e.URL)
	duration := time.Since(start)
	if err != nil {
		log.Printf("External API health check failed: %v", err)
		return entities.HealthCheckResult{Status: entities.HealthStatusDown, Error: err, Duration: duration.String()}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("External API health check failed: received status code %d", resp.StatusCode)
		return entities.HealthCheckResult{Status: entities.HealthStatusDown, Error: fmt.Errorf("received status code %d", resp.StatusCode), Duration: duration.String()}
	}
	log.Printf("External API health check successful: %s", duration.String())
	return entities.HealthCheckResult{Status: entities.HealthStatusUp, Duration: duration.String()}
}

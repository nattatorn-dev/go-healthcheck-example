package entities

type HealthStatus string

const (
	HealthStatusUp   HealthStatus = "UP"
	HealthStatusDown HealthStatus = "DOWN"
)

type HealthCheckResult struct {
	Status   HealthStatus `json:"status"`
	Error    error        `json:"error,omitempty"` // Detailed error information
	Duration string       `json:"duration"`
}

type HealthStatusEntry struct {
	Name     string       `json:"name"`
	Status   HealthStatus `json:"status"`
	Error    string       `json:"error,omitempty"` // Store the error message as a string
	Duration string       `json:"duration"`
}

type HealthCheckResponse struct {
	Status HealthStatus        `json:"status"`
	Checks []HealthStatusEntry `json:"checks"`
}

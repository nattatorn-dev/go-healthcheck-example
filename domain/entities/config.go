package entities

import "time"

type HealthCheckConfig struct {
	Timeout  time.Duration
	Interval time.Duration
}

type CheckerConfig struct {
	Timeout  time.Duration
	Interval time.Duration
}

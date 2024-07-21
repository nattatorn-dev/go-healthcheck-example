package services

import (
	"log"
	"time"

	"github.com/nattatorn-dev/go-healthcheck/domain/entities"
	"github.com/nattatorn-dev/go-healthcheck/domain/repositories"
)

type HealthService struct {
	readinessCheckers map[string]checkerConfigPair
	livenessCheckers  map[string]checkerConfigPair
	store             *entities.HealthStore
}

type checkerConfigPair struct {
	Checker repositories.HealthChecker
	Config  entities.CheckerConfig
}

func NewHealthService(store *entities.HealthStore) *HealthService {
	return &HealthService{
		readinessCheckers: make(map[string]checkerConfigPair),
		livenessCheckers:  make(map[string]checkerConfigPair),
		store:             store,
	}
}

func (s *HealthService) RegisterReadiness(name string, checker repositories.HealthChecker, config entities.CheckerConfig) {
	s.readinessCheckers[name] = checkerConfigPair{Checker: checker, Config: config}
}

func (s *HealthService) RegisterLiveness(name string, checker repositories.HealthChecker, config entities.CheckerConfig) {
	s.livenessCheckers[name] = checkerConfigPair{Checker: checker, Config: config}
}

func (s *HealthService) CheckAllHealth() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	for name, pair := range s.readinessCheckers {
		status := pair.Checker.CheckHealth()
		s.store.SetStatus(name+"-readiness", entities.HealthStatusEntry{
			Name:     name,
			Status:   status.Status,
			Error:    getErrorMessage(status.Error),
			Duration: status.Duration,
		})
	}
	for name, pair := range s.livenessCheckers {
		status := pair.Checker.CheckHealth()
		s.store.SetStatus(name+"-liveness", entities.HealthStatusEntry{
			Name:     name,
			Status:   status.Status,
			Error:    getErrorMessage(status.Error),
			Duration: status.Duration,
		})
	}
}

func (s *HealthService) GetReadinessStatuses() []entities.HealthStatusEntry {
	return s.store.GetStatusesByType("readiness")
}

func (s *HealthService) GetLivenessStatuses() []entities.HealthStatusEntry {
	return s.store.GetStatusesByType("liveness")
}

func (s *HealthService) startChecker(name string, pair checkerConfigPair) {
	ticker := time.NewTicker(pair.Config.Interval)
	go func() {
		for {
			<-ticker.C
			log.Printf("Performing scheduled health check for %s", name)
			status := pair.Checker.CheckHealth()
			s.store.SetStatus(name+"-liveness", entities.HealthStatusEntry{
				Name:     name,
				Status:   status.Status,
				Error:    getErrorMessage(status.Error),
				Duration: status.Duration,
			})
		}
	}()
}

func (s *HealthService) StartBackgroundCheck() {
	s.CheckAllHealth() // Perform an initial check immediately
	for name, pair := range s.readinessCheckers {
		s.startChecker(name, pair)
	}
	for name, pair := range s.livenessCheckers {
		s.startChecker(name, pair)
	}
}

func getErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

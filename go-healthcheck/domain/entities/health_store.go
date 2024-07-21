package entities

import "sync"

type HealthStore struct {
	statuses map[string]HealthStatusEntry
	mutex    sync.RWMutex
}

func NewHealthStore() *HealthStore {
	return &HealthStore{
		statuses: make(map[string]HealthStatusEntry),
	}
}

func (s *HealthStore) SetStatus(name string, status HealthStatusEntry) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.statuses[name] = status
}

func (s *HealthStore) GetStatusesByType(checkType string) []HealthStatusEntry {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var statuses []HealthStatusEntry
	for name, status := range s.statuses {
		if checkType == "readiness" && len(name) > 9 && name[len(name)-9:] == "readiness" {
			statuses = append(statuses, status)
		}
		if checkType == "liveness" && len(name) > 8 && name[len(name)-8:] == "liveness" {
			statuses = append(statuses, status)
		}
	}
	return statuses
}

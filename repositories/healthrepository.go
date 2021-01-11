package repositories

import (
	"sync"
	"time"

	"github.com/jake-hansen/followrs/domain"
)

type startTime struct {
	mu        sync.Mutex
	startTime time.Time
}

// HealthRepository represents a repository for retrieving Health.
type HealthRepository struct {
	timeInfo startTime
}

// NewSimpleHealthRepository will create an implementation of repositories.HealthRepository.
func NewSimpleHealthRepository(startTime time.Time) domain.HealthRepository {
	healthRepo := &HealthRepository{}
	healthRepo.timeInfo.setStartTime(startTime)
	return healthRepo
}

// GetStatus returns the string "online"
func (hr *HealthRepository) GetStatus() (string, error) {
	return "online", nil
}

// GetUpsince returns the time the application started.
func (hr *HealthRepository) GetUpsince() (time.Time, error) {
	return hr.timeInfo.startTime, nil
}

// SetStartTime sets the start time of the server to the provided time.
func (hr *HealthRepository) SetStartTime(time time.Time) {
	hr.timeInfo.setStartTime(time)
}

func (t *startTime) setStartTime(time time.Time) {
	t.mu.Lock()
	t.startTime = time
	t.mu.Unlock()
}

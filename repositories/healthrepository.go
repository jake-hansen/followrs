package repositories

import (
	"time"

	"github.com/jake-hansen/followrs/domain"
)

// HealthRepository represents a repository for retrieving Health.
type HealthRepository struct{}

// NewSimpleHealthRepository will create an implementation of repositories.HealthRepository.
func NewSimpleHealthRepository() domain.HealthRepository {
	return &HealthRepository{}
}

// GetStatus returns the string "online"
func (hr *HealthRepository) GetStatus() (string, error) {
	return "online", nil
}

// GetUpsince is not implemented.
func (hr *HealthRepository) GetUpsince() (time.Time, error) {
	panic("not implemented") // TODO: Implement
}

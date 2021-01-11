package repositories

import (
	"github.com/jake-hansen/followrs/domain"
)

type HealthRepository struct{}

func NewSimpleHealthRepository() domain.HealthRepository {
	return &HealthRepository{}
}

func (hr *HealthRepository) GetStatus() (string, error) {
	return "online", nil
}

func (hr *HealthRepository) GetUpsince() (string, error) {
	panic("not implemented") // TODO: Implement
}

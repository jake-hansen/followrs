package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// HealthRepository is a mock HealthRepository.
type HealthRepository struct {
	mock.Mock
}

// GetStatus provides a mock function.
func (m *HealthRepository) GetStatus() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// GetUpsince provides a mock function.
func (m *HealthRepository) GetUpsince() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

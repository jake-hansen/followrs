package mocks

import (
	"github.com/jake-hansen/followrs/domain"
	"github.com/stretchr/testify/mock"
)

type HealthService struct {
	mock.Mock
}

func (m *HealthService) GetHealth() (domain.Health, error) {
	args := m.Called()
	return args.Get(0).(domain.Health), args.Error(1)
}

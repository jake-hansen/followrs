package services

import (
	"github.com/jake-hansen/followrs/domain"
)

type HealthService struct {
	Repo domain.HealthRepository
}

func NewSimpleHealthService(repo domain.HealthRepository) domain.HealthService {
	return &HealthService{
		Repo: repo,
	}
}

func (hs *HealthService) GetHealth() (domain.Health, error) {
	status, _ := hs.Repo.GetStatus()
	return domain.Health{
		Status: status,
	}, nil
}

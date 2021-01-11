package services

import (
	"errors"

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
	status, err := hs.Repo.GetStatus()
	var returnErr error = nil
	if err != nil {
		returnErr = errors.New("an error occurred retrieving health status of the server")
	}
	return domain.Health{
		Status: status,
	}, returnErr
}

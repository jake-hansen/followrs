package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jake-hansen/followrs/repositories/mocks"
	healthservice "github.com/jake-hansen/followrs/services"
)

// TestGetHealth tests HealthService's GetHealth func.
func TestGetHealth(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.HealthRepository)
		statusStr := "online"
		repo.On("GetStatus").Return(statusStr, nil)
		service := healthservice.NewSimpleHealthService(repo)

		health, err := service.GetHealth()

		assert.Equal(t, health.Status, statusStr)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		repo := new(mocks.HealthRepository)
		repo.On("GetStatus").Return("", errors.New("example error"))
		service := healthservice.NewSimpleHealthService(repo)
		health, err := service.GetHealth()

		assert.Equal(t, health.Status, "")
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

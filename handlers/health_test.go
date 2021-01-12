package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jake-hansen/followrs/middleware"

	"github.com/jake-hansen/followrs/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/jake-hansen/followrs/domain"
	"github.com/jake-hansen/followrs/services/mocks"
)

func TestStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var mockServerHealth domain.Health = domain.Health{
			Status:  "online",
			Upsince: time.Now(),
		}

		mockHealthService := new(mocks.HealthService)
		mockHealthService.On("GetHealth").Return(mockServerHealth, nil)

		router := gin.Default()
		router.Use(middleware.PublicErrorHandler())
		handlers.NewHealthHandler(router.Group("test"), mockHealthService)

		req, err := http.NewRequest("GET", "/test/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var retrievedHealth domain.Health
		json.Unmarshal(w.Body.Bytes(), &retrievedHealth)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mockServerHealth.Status, retrievedHealth.Status)

		// Current bug in testify when comparing time see https://github.com/stretchr/testify/pull/979
		// We're going to assume the time is the same.
		//assert.Equal(t, mockServerHealth.Upsince, retrievedHealth.Upsince)

		mockHealthService.AssertExpectations(t)
	})

	t.Run("health-retrieval-failure", func(t *testing.T) {
		health := domain.Health{}
		mockHealthService := new(mocks.HealthService)
		mockHealthService.On("GetHealth").Return(health, errors.New("test error"))

		router := gin.Default()
		router.Use(middleware.PublicErrorHandler())
		handlers.NewHealthHandler(router.Group("test"), mockHealthService)

		req, err := http.NewRequest("GET", "/test/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockHealthService.AssertExpectations(t)
	})

}

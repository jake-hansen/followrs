package domain

import (
	"time"
)

// Health represents health of the server.
type Health struct {
	Status string `json:"status"`
}

type HealthService interface {
	GetHealth() (Health, error)
}

type HealthRepository interface {
	GetStatus() (string, error)
	GetUpsince() (time.Time, error)
}

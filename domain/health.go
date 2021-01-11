package domain

import (
	"time"
)

// Health represents health of the server.
type Health struct {
	Status  string    `json:"status"`
	Upsince time.Time `json:"upsince"`
}

type HealthService interface {
	GetHealth() (Health, error)
}

type HealthRepository interface {
	GetStatus() (string, error)
	GetUpsince() (time.Time, error)
	SetStartTime(time time.Time)
}

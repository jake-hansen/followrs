package repositories_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jake-hansen/followrs/repositories"
)

func TestGetStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := repositories.NewSimpleHealthRepository(time.Now())

		status, err := repo.GetStatus()

		assert.NoError(t, err)
		assert.Equal(t, "online", status)
	})
}

func TestGetUpsince(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Now()
		repo := repositories.NewSimpleHealthRepository(now)

		time, err := repo.GetUpsince()

		assert.NoError(t, err)
		assert.Equal(t, now, time)
	})
}

func TestSetStartTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Now()
		repo := repositories.NewSimpleHealthRepository(now)

		d, _ := time.ParseDuration("1m")

		repo.SetStartTime(now.Add(d))
		time, err := repo.GetUpsince()

		assert.NoError(t, err)
		assert.Equal(t, now.Add(d), time)
	})
}

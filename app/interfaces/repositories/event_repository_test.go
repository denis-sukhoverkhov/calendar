package repositories

import (
	"github.com/denis-sukhoverkhov/calendar/app/domain"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventRepository(t *testing.T) {
	now := time.Now()

	t.Run("Store event", func(t *testing.T) {
		eventRepo := NewEventRepository()
		from, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T10:00:00+00:00")
		to, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T15:00:00+00:00")
		event := models.Event{Id: 1, Name: "work", From: from, To: to, UserId: 2, Active: true, CreatedAt: &now, UpdatedAt: nil}
		eventRepo.Store(event)

		from2, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T11:00:00+00:00")
		to2, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T12:00:00+00:00")

		event2 := models.Event{Id: 2, Name: "birthday", From: from2, To: to2, UserId: 2, Active: true, CreatedAt: &now, UpdatedAt: nil}
		_, err := eventRepo.Store(event2)
		assert.Equal(t, err, domain.ErrDateBusy)

	})

	t.Run("Store event", func(t *testing.T) {
		eventRepo := NewEventRepository()
		from, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T10:00:00+00:00")
		to, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T15:00:00+00:00")
		event := models.Event{Id: 1, Name: "work", From: from, To: to, UserId: 2, Active: true, CreatedAt: &now, UpdatedAt: nil}
		_, _ = eventRepo.Store(event)

		from2, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T15:00:00+00:00")
		to2, _ := time.Parse(
			time.RFC3339,
			"2020-06-06T17:00:00+00:00")

		event2 := models.Event{Id: 2, Name: "birthday", From: from2, To: to2, UserId: 2, Active: true, CreatedAt: &now, UpdatedAt: nil}
		storedEvent, err := eventRepo.Store(event2)
		assert.Equal(t, err, nil)

		expected := storedEvent
		actual := eventRepo.FindById(event2.Id)
		assert.Equal(t, expected, actual)
	})

}

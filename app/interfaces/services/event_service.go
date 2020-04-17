package services

import (
	"fmt"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
)

type eventService struct {
	rep repositories.EventRepository
}

func NewEventService(repository *repositories.EventRepository) *eventService {
	return &eventService{
		rep: *repository,
	}
}

func (e *eventService) GetById(eventId int) (*models.Event, error) {
	event, err := e.rep.FindById(eventId)
	if err != nil {
		return nil, fmt.Errorf("EventService.GetById error %w", err)
	}
	return event, nil
}

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

func (e *eventService) GetAll() ([]*models.Event, error) {
	event, err := e.rep.FindAll()
	if err != nil {
		return nil, fmt.Errorf("EventService.GetAll error %w", err)
	}
	return event, nil
}

func (e *eventService) Delete(eventId int) error {
	err := e.rep.Delete(eventId)
	if err != nil {
		return fmt.Errorf("EventService.Delete error %w", err)
	}
	return nil
}

package services

import (
	"fmt"
	"github.com/denis-sukhoverkhov/calendar/internal/domain"
	"github.com/denis-sukhoverkhov/calendar/internal/domain/models"
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/repositories"
	"time"
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

func (e *eventService) Save(event *models.Event) (*models.Event, error) {
	alreadyStoredEventsForCurrentUser, err := e.rep.FindByUserId(event.UserId)
	if len(alreadyStoredEventsForCurrentUser) > 0 {
		for i := 0; i < len(alreadyStoredEventsForCurrentUser); i++ {
			if (event.From.After(alreadyStoredEventsForCurrentUser[i].From) &&
				event.From.Before(alreadyStoredEventsForCurrentUser[i].To)) ||
				(event.To.After(alreadyStoredEventsForCurrentUser[i].From) &&
					event.From.Before(alreadyStoredEventsForCurrentUser[i].To)) ||
				(event.From.Before(alreadyStoredEventsForCurrentUser[i].From) && event.To.After(alreadyStoredEventsForCurrentUser[i].To)) {
				return nil, domain.ErrDateBusy
			}
		}
	}

	newEvent, err := e.rep.Store(*event)
	if err != nil {
		return nil, fmt.Errorf("EventService.Save error %w", err)
	}
	return newEvent, nil
}

func (e *eventService) Delete(eventId int) error {
	err := e.rep.Delete(eventId)
	if err != nil {
		return fmt.Errorf("EventService.Delete error %w", err)
	}
	return nil
}

func (e *eventService) GetAllByDay(userId int64, date time.Time) ([]*models.Event, error) {
	events, err := e.rep.FindByUserIdAndDate(userId, date)
	if err != nil {
		return nil, fmt.Errorf("EventService.GetAllByDay error %w", err)
	}
	return events, nil
}

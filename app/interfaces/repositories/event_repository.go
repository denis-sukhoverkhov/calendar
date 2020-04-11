package repositories

import (
	"errors"
	"github.com/denis-sukhoverkhov/calendar/app/domain"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
)

type EventRepository interface {
	Store(user models.Event)
	FindById(id int)
	FindAll()
	Delete(id int)
}

type eventRepository struct {
	events map[int64]*models.Event
}

func NewEventRepository() *eventRepository {
	return &eventRepository{
		events: map[int64]*models.Event{},
	}
}

func (r *eventRepository) FindById(id int64) *models.Event {
	return r.events[id]
}

func (r *eventRepository) FindByUserId(Userid int64) []models.Event {
	events := make([]models.Event, 0)
	for _, val := range r.events {
		if val.UserId == Userid {
			events = append(events, *val)
		}
	}
	return events
}

func (r *eventRepository) FindAll() []models.Event {

	events := make([]models.Event, len(r.events))
	for _, val := range r.events {
		events[val.Id-1] = *val
	}
	return events
}

func (r *eventRepository) Store(event models.Event) (*models.Event, error) {
	alreadyStoredEventsForCurrentUser := r.FindByUserId(event.UserId)
	if len(alreadyStoredEventsForCurrentUser) == 0 {
		r.events[event.Id] = &event
		return r.events[event.Id], nil
	}

	for i := 0; i < len(alreadyStoredEventsForCurrentUser); i++ {

		if (event.From.After(alreadyStoredEventsForCurrentUser[i].From) &&
			event.From.Before(alreadyStoredEventsForCurrentUser[i].To)) ||
			(event.To.After(alreadyStoredEventsForCurrentUser[i].From) &&
				event.From.Before(alreadyStoredEventsForCurrentUser[i].To)) ||
			(event.From.Before(alreadyStoredEventsForCurrentUser[i].From) && event.To.After(alreadyStoredEventsForCurrentUser[i].To)) {
			return nil, domain.ErrDateBusy
		}
	}
	r.events[event.Id] = &event
	return r.events[event.Id], nil
}

func (r *eventRepository) Delete(id int64) error {
	if _, ok := r.events[id]; ok {
		delete(r.events, id)
		return nil
	}

	return errors.New("removing user does not exist in userRepository")
}

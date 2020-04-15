package db

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userDbRepository struct {
	pool *pgxpool.Pool
}

func NewUserDbRepository(pgPool *pgxpool.Pool) repositories.UserRepository {
	return &userDbRepository{
		pool: pgPool,
	}
}

func (r *userDbRepository) FindById(id int) (*models.User, error) {
	query := sq.Select("*").From("\"user\"").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("User.FindById QueryBuilder error %w", err)
	}

	user := &models.User{}
	err = r.pool.QueryRow(context.Background(), sql, args...).Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("User.FindById QueryRow error %w", err)
	}
	return user, nil
}

//func (r *eventRepository) FindByUserId(Userid int64) []models.Event {
//	events := make([]models.Event, 0)
//	for _, val := range r.events {
//		if val.UserId == Userid {
//			events = append(events, *val)
//		}
//	}
//	return events
//}
//
//func (r *eventRepository) FindAll() []models.Event {
//
//	events := make([]models.Event, len(r.events))
//	for _, val := range r.events {
//		events[val.Id-1] = *val
//	}
//	return events
//}
//
//func (r *eventRepository) Store(event models.Event) (*models.Event, error) {
//	alreadyStoredEventsForCurrentUser := r.FindByUserId(event.UserId)
//	if len(alreadyStoredEventsForCurrentUser) == 0 {
//		r.events[event.Id] = &event
//		return r.events[event.Id], nil
//	}
//
//	for i := 0; i < len(alreadyStoredEventsForCurrentUser); i++ {
//
//		if (event.From.After(alreadyStoredEventsForCurrentUser[i].From) &&
//			event.From.Before(alreadyStoredEventsForCurrentUser[i].To)) ||
//			(event.To.After(alreadyStoredEventsForCurrentUser[i].From) &&
//				event.From.Before(alreadyStoredEventsForCurrentUser[i].To)) ||
//			(event.From.Before(alreadyStoredEventsForCurrentUser[i].From) && event.To.After(alreadyStoredEventsForCurrentUser[i].To)) {
//			return nil, domain.ErrDateBusy
//		}
//	}
//	r.events[event.Id] = &event
//	return r.events[event.Id], nil
//}
//
//func (r *eventRepository) Delete(id int64) error {
//	if _, ok := r.events[id]; ok {
//		delete(r.events, id)
//		return nil
//	}
//
//	return errors.New("removing user does not exist in userRepository")
//}

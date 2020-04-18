package db

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type evenDbRepository struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewEventDbRepository(pgPool *pgxpool.Pool) repositories.EventRepository {
	return &evenDbRepository{
		pool: pgPool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *evenDbRepository) FindById(id int) (*models.Event, error) {
	query := r.sq.Select("*").From("\"event\"").Where(sq.Eq{"id": id}, sq.Eq{"active": true})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("Event.FindById QueryBuilder error %w", err)
	}

	event := &models.Event{}
	err = r.pool.QueryRow(context.Background(), sql, args...).Scan(
		&event.Id, &event.Name, &event.From, &event.To, &event.UserId, &event.Active, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Event.FindById QueryRow error %w", err)
	}
	return event, nil
}

func (r *evenDbRepository) FindAll() ([]*models.Event, error) {
	query := r.sq.Select("*").From("\"event\"").Where(sq.Eq{"active": true})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("Event.FindAll QueryBuilder error %w", err)
	}

	events := make([]*models.Event, 0)
	rows, err := r.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Event.FindAll rows errors %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		event := models.Event{}
		err := rows.Scan(
			&event.Id, &event.Name, &event.From, &event.To, &event.UserId, &event.Active, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("Event.FindAll QueryRow error %w", err)
		}
		events = append(events, &event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("Event.FindAll all errors %w", err)
	}

	return events, nil
}

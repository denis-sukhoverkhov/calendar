package interfaces

import (
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/repositories"
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/repositories/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitRepositories(pgPool *pgxpool.Pool) *repositories.RepositoryInteractor {

	return &repositories.RepositoryInteractor{
		User:  db.NewUserDbRepository(pgPool),
		Event: db.NewEventDbRepository(pgPool),
	}
}

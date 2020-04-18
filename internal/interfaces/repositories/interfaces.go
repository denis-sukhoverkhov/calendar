package repositories

import (
	"github.com/denis-sukhoverkhov/calendar/internal/domain/models"
)

type EventRepository interface {
	Store(user models.Event) (*models.Event, error)
	FindById(id int) (*models.Event, error)
	FindByUserId(userId int64) ([]*models.Event, error)
	FindAll() ([]*models.Event, error)
	Delete(id int) error
}

type UserRepository interface {
	Store(user models.User) (*models.User, error)
	FindById(id int) (*models.User, error)
	FindAll() ([]*models.User, error)
	Delete(id int) error
}

type RepositoryInteractor struct {
	User  UserRepository
	Event EventRepository
}

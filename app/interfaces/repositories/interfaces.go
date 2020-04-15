package repositories

import (
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
)

type EventRepository interface {
	Store(user models.Event)
	FindById(id int)
	FindAll()
	Delete(id int)
}

type UserRepository interface {
	//Store(user models.User) *models.User
	FindById(id int) (*models.User, error)
	//FindAll() []*models.User
	//Delete(id int) error
}

type RepositoryInteractor struct {
	User  UserRepository
	Event EventRepository
}

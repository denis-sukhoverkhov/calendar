package interfaces

import (
	"calendar/calendar/domain/models"
	"errors"
)

type UserRepository interface {
	Store(user models.User)
	FindById(id int)
	FindAll()
	Delete(id int)
}

type userRepository struct {
	users map[int64]*models.User
}

func NewUserRepository() *userRepository {
	return &userRepository{
		users: map[int64]*models.User{},
	}
}

func (r *userRepository) FindById(id int64) *models.User {
	return r.users[id]
}

func (r *userRepository) FindAll() []models.User {

	users := make([]models.User, len(r.users))
	for _, val := range r.users {
		users[val.Id-1] = *val
	}
	return users
}

func (r *userRepository) Store(user models.User) models.User {
	r.users[user.Id] = &user
	return *r.users[user.Id]
}

func (r *userRepository) Delete(id int64) error {
	if _, ok := r.users[id]; ok {
		delete(r.users, id)
		return nil
	}

	return errors.New("removing user does not exist in userRepository")
}

type EventRepository interface {
	Store(event models.Event)
	FindById(id int)
}

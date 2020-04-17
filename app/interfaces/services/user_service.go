package services

import (
	"fmt"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
)

type userService struct {
	rep repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *userService {
	return &userService{
		rep: *repository,
	}
}

func (u *userService) GetById(userId int) (*models.User, error) {
	user, err := u.rep.FindById(userId)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetById error %w", err)
	}
	return user, nil
}

func (u *userService) GetAll() ([]*models.User, error) {
	user, err := u.rep.FindAll()
	if err != nil {
		return nil, fmt.Errorf("UserService.GetAll error %w", err)
	}
	return user, nil
}

func (u *userService) Save(user *models.User) (*models.User, error) {
	newUser, err := u.rep.Store(*user)
	if err != nil {
		return nil, fmt.Errorf("UserService.Save error %w", err)
	}
	return newUser, nil
}

func (u *userService) Delete(userId int) error {
	err := u.rep.Delete(userId)
	if err != nil {
		return fmt.Errorf("UserService.Delete error %w", err)
	}
	return nil
}

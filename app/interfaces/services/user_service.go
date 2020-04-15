package services

import (
	"fmt"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
)

type UserService struct {
	rep repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{
		rep: *repository,
	}
}

func (u *UserService) GetById(userId int) (*models.User, error) {
	user, err := u.rep.FindById(userId)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetById error %w", err)
	}
	return user, nil
}

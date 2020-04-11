package repositories

import (
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserRepository(t *testing.T) {
	now := time.Now()

	t.Run("Store to userRepo", func(t *testing.T) {
		userRepo := NewUserRepository()
		user := models.User{Id: 1, FirstName: "Denis", LastName: "Sukhoverkhov", Active: true, CreatedAt: &now, UpdatedAt: nil}
		userRepo.Store(user)
		expected := &user
		actual := userRepo.FindById(user.Id)
		assert.Equal(t, expected, actual)
	})

	t.Run("FindAll users", func(t *testing.T) {
		userRepo := NewUserRepository()
		user1 := models.User{Id: 1, FirstName: "Denis", LastName: "Sukhoverkhov", Active: true, CreatedAt: &now, UpdatedAt: nil}
		user2 := models.User{Id: 2, FirstName: "Denis", LastName: "Sukhoverkhov", Active: true, CreatedAt: &now, UpdatedAt: nil}

		userRepo.Store(user1)
		userRepo.Store(user2)

		expected := []models.User{
			user1,
			user2,
		}
		actual := userRepo.FindAll()
		assert.Equal(t, expected, actual)
	})

	t.Run("Delete user from UserRepository", func(t *testing.T) {
		userRepo := NewUserRepository()
		user1 := models.User{Id: 1, FirstName: "Denis", LastName: "Sukhoverkhov", Active: true, CreatedAt: &now, UpdatedAt: nil}
		user2 := models.User{Id: 2, FirstName: "Denis", LastName: "Sukhoverkhov", Active: true, CreatedAt: &now, UpdatedAt: nil}

		userRepo.Store(user1)
		userRepo.Store(user2)

		err := userRepo.Delete(user2.Id)
		if err != nil {
			assert.Errorf(t, err, "Ошибка удаления пользователя из репозитория")
		}

		expected := []models.User{
			user1,
		}
		actual := userRepo.FindAll()
		assert.Equal(t, expected, actual)
	})

}

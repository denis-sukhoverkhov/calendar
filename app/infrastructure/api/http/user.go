package infrastructure

import (
	"encoding/json"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/services"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func GetUserHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(chi.URLParam(r, "userId"))
		userService := services.NewUserService(&repos.User)
		user, err := userService.GetById(userId)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(user)
		writer.Write(userJson)
	}
}

func GetUsersHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		userService := services.NewUserService(&repos.User)
		users, err := userService.GetAll()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(users)
		writer.Write(userJson)
	}
}

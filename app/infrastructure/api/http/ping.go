package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/services"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
)

func Hello(writer http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Hello world!"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func Main(writer http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Welcome to the main page!"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func NewHelloHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
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

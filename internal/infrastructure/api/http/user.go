package infrastructure

import (
	"encoding/json"
	"errors"
	"github.com/denis-sukhoverkhov/calendar/internal/domain/models"
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/repositories"
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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

func PostUserHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		data := &UserRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(writer, r, ErrInvalidRequest(err))
			return
		}
		userService := services.NewUserService(&repos.User)
		user, err := userService.Save(&models.User{FirstName: data.FirstName, LastName: data.LastName})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(user)
		writer.Write(userJson)
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a *UserRequest) Bind(r *http.Request) error {
	if a.FirstName == "" {
		return errors.New("missing required FirstName fields")
	}
	if a.LastName == "" {
		return errors.New("missing required LastName fields")
	}
	return nil
}

func DeleteUserHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(chi.URLParam(r, "userId"))
		userService := services.NewUserService(&repos.User)
		err := userService.Delete(userId)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
		}
		writer.WriteHeader(http.StatusOK)
	}
}

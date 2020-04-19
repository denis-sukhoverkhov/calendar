package infrastructure

import (
	"encoding/json"
	"errors"
	"github.com/denis-sukhoverkhov/calendar/internal/domain"
	"github.com/denis-sukhoverkhov/calendar/internal/domain/models"
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/repositories"
	"github.com/denis-sukhoverkhov/calendar/internal/interfaces/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"time"
)

func GetVentHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		eventId, _ := strconv.Atoi(chi.URLParam(r, "eventId"))
		eventService := services.NewEventService(&repos.Event)
		event, err := eventService.GetById(eventId)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			//writer.Write([]byte("500 - Something bad happened!"))
			return
		}
		if event == nil {
			writer.WriteHeader(http.StatusNotFound)
			//writer.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(event)
		writer.Write(userJson)
	}
}

func GetEventsHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		eventService := services.NewEventService(&repos.Event)
		events, err := eventService.GetAll()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(events)
		writer.Write(userJson)
	}
}

func PostEventHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		data := &EventRequest{}
		if err := render.Bind(r, data); err != nil {
			_ = render.Render(writer, r, ErrInvalidRequest(err))
			return
		}
		eventService := services.NewEventService(&repos.Event)
		user, err := eventService.Save(&models.Event{Name: data.Name, From: data.From, To: data.To, UserId: data.UserId})
		if err != nil {
			if err == domain.ErrDateBusy {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(err.Error()))
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(user)
		writer.Write(userJson)
	}
}

type EventRequest struct {
	Name   string    `json:"name"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	UserId int64     `json:"user_id"`
}

func (a *EventRequest) Bind(r *http.Request) error {
	if a.Name == "" {
		return errors.New("missing required Name fields")
	}
	if &a.From == nil {
		return errors.New("missing required From fields")
	}
	if &a.To == nil {
		return errors.New("missing required To fields")
	}
	if a.UserId == 0 {
		return errors.New("missing required UserId fields")
	}
	return nil
}

func DeleteEventHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		eventId, _ := strconv.Atoi(chi.URLParam(r, "eventId"))
		eventService := services.NewEventService(&repos.Event)
		err := eventService.Delete(eventId)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func GetEventsForDayHandler(repos *repositories.RepositoryInteractor) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(r.URL.Query().Get("userId"))
		date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
			return
		}

		eventService := services.NewEventService(&repos.Event)
		events, err := eventService.GetAllByDay(int64(userId), date)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 - Something bad happened!"))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		userJson, err := json.Marshal(events)
		writer.Write(userJson)
	}
}

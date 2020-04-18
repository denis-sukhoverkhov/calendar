package infrastructure

import (
	"context"
	"fmt"
	apihttp "github.com/denis-sukhoverkhov/calendar/app/infrastructure/api/http"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppServer struct {
	l    *zap.Logger
	repo *repositories.RepositoryInteractor
	*http.Server
}

func NewServer(config Configuration) (*AppServer, error) {

	logger := NewLogger(config.Logs.Level, config.Logs.PathToLogFile)
	listenAddr := fmt.Sprintf("%s:%d", config.HttpListen.Ip, config.HttpListen.Port)
	// FIXME: эта строка нужна?
	errorLog, _ := zap.NewStdLogAt(logger, zap.ErrorLevel)

	pgPool := NewPgPool(config, logger)
	repos := interfaces.InitRepositories(pgPool)

	srv := http.Server{
		Addr:         listenAddr,
		Handler:      NewHttpApi(repos, logger),
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &AppServer{logger, repos, &srv}, nil
}

func (s *AppServer) Start() {
	s.l.Info("Starting server")
	defer s.l.Sync()

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.l.Fatal("Could not listen on", zap.String("addr", s.Addr), zap.Error(err))
		}
	}()

	s.l.Info("Server is ready to handle requests", zap.String("addr", s.Addr))
	s.gracefulShutdown()
}

func (s *AppServer) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	s.l.Info("Server is shutting down", zap.String("reason", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		s.l.Fatal("Could not gracefully shutdown the server", zap.Error(err))
	}
	s.l.Info("Server stopped")
}

func NewHttpApi(repos *repositories.RepositoryInteractor, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(zapLogger(logger))
	r.Use(middleware.Recoverer)

	r.Get("/", apihttp.Main)
	r.Get("/hello", apihttp.Hello)

	r.Get("/user/{userId:[0-9]+}", apihttp.GetUserHandler(repos))
	r.Get("/user", apihttp.GetUsersHandler(repos))
	r.Post("/user", apihttp.PostUserHandler(repos))
	r.Delete("/user/{userId:[0-9]+}", apihttp.DeleteUserHandler(repos))

	r.Get("/event/{eventId:[0-9]+}", apihttp.GetVentHandler(repos))
	r.Get("/event", apihttp.GetEventsHandler(repos))
	r.Delete("/event/{eventId:[0-9]+}", apihttp.DeleteEventHandler(repos))

	logRoutes(r, logger)
	return r
}

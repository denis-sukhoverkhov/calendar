package calendar

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	l *zap.Logger
	*http.Server
}

func NewServer(listenAddr string, logLevel zapcore.LevelEnabler, mux http.Handler) (*Server, error) {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(logLevel))
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	errorLog, _ := zap.NewStdLogAt(logger, zap.ErrorLevel)
	srv := http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{logger, &srv}, nil
}

func (s *Server) Start() {
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

func (s *Server) gracefulShutdown() {
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

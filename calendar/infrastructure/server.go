package infrastructure

import (
	api_http "calendar/calendar/infrastructure/api/http"
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	l *zap.Logger
	*http.Server
}

func NewServer(listenAddr string, logLevel zapcore.LevelEnabler) (*Server, error) {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(logLevel))
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	errorLog, _ := zap.NewStdLogAt(logger, zap.ErrorLevel)
	srv := http.Server{
		Addr:         listenAddr,
		Handler:      NewHttpApi(logger),
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{logger, &srv}, nil
}

func NewHttpApi(logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(zapLogger(logger))
	r.Use(middleware.Recoverer)

	r.Handle("/hello", http.HandlerFunc(api_http.Hello))
	r.Handle("/", http.HandlerFunc(api_http.HandleNotFound))

	logRoutes(r, logger)

	return r
}

func zapLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("Served",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func logRoutes(r *chi.Mux, logger *zap.Logger) {
	if err := chi.Walk(r, zapPrintRoute(logger)); err != nil {
		logger.Error("Failed to walk routes:", zap.Error(err))
	}
}

func zapPrintRoute(logger *zap.Logger) chi.WalkFunc {
	return func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		logger.Debug("Registering route", zap.String("method", method), zap.String("route", route))
		return nil
	}
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

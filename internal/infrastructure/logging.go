package infrastructure

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

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

func NewLogger(level string, pathToLogFile string) *zap.Logger {

	err := os.MkdirAll(path.Dir(pathToLogFile), os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory for logs, %v", err)
	}

	logLevel := zap.NewAtomicLevel()
	err = logLevel.UnmarshalText([]byte(level))
	if err != nil {
		log.Fatalf("Unmarshaling error of logger config, %v", err)
	}

	cfg := zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout", "stderr", pathToLogFile},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build(zap.AddStacktrace(logLevel))
	if err != nil {
		log.Fatalf("Logger error creating, %v", err)
	}

	logger.Info("logger construction succeeded")

	return logger
}

package calendar

import (
	"go.uber.org/zap"
	"net/http"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, _ := zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
		sugar := logger.Sugar()
		next.ServeHTTP(w, r)
		sugar.Infow("fetch URL",
			"url", r.URL.Path,
			"method", r.Method,
		)
	})
}

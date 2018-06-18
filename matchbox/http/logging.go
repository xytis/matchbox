package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *loggingResponseWriter) Write(buf []byte) (int, error) {
	size, err := w.ResponseWriter.Write(buf)
	// Write triggers WriteHeader
	if w.status == 0 && err == nil {
		w.status = http.StatusOK
	}
	w.size = size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.status = statusCode
}

// logging logs HTTP requests.
func (s *Server) logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		lw := &loggingResponseWriter{w, 0, 0}

		ctx := req.Context()
		var rID string
		if rID = req.Header.Get("X-Request-ID"); rID != "" {
		} else {
			r, err := uuid.NewRandom()
			if err != nil {
				http.Error(w, fmt.Sprintf(`meteorite collided with random generator function: %v`, err), http.StatusInternalServerError)
				return
			}
			rID = r.String()
		}
		ctx = withRequestID(ctx, rID)

		next.ServeHTTP(lw, req.WithContext(ctx))

		duration := time.Since(startTime)
		s.logger.Info("HTTP",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Duration("duration", duration),
			zap.Int("code", lw.status),
			zap.Int("size", lw.size),
			zap.String("request-id", rID),
		)
	}
	return http.HandlerFunc(fn)
}

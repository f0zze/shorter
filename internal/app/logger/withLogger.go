package logger

import (
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func NewLogger(fileName string) zerolog.Logger {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	logger := zerolog.New(file).With().Timestamp().Logger()

	return logger
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.responseData.status = status
}

func WithLogging(l *zerolog.Logger) func(h http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(h http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			var responseData = &responseData{}
			responseWriter := loggingResponseWriter{
				w, responseData,
			}

			h.ServeHTTP(&responseWriter, r)
			l.Info().
				Str("uri", r.RequestURI).
				Str("method", r.Method).
				Str("duration", time.Since(startTime).String()).
				Msg("New request")
			l.Info().
				Int("status", responseData.status).
				Int("size", responseData.size).
				Msg("Request response")
		}
	}

}

package middlewares

import (
	"net/http"

	"gitlab.com/hydra/forum-api/api/logger"
)

// LoggingResponseWriter is a wrapper around an http.ResponseWriter which captures the
// status code written to the response, so that it can be logged.
type LoggingResponseWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
	// Response content could also be captured here, but I was only interested in logging the response status code
}

func NewLoggingResponseWriter(wrapped http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{wrapped: wrapped}
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.wrapped.Header()
}

func (lrw *LoggingResponseWriter) Write(content []byte) (int, error) {
	return lrw.wrapped.Write(content)
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.wrapped.WriteHeader(statusCode)
}

// Log is a middleware method for writing log for every passed request
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging stuff
		w2 := NewLoggingResponseWriter(w)
		next.ServeHTTP(w2, r)

		log := logger.GetInstance()
		log.Infof("%s %s, response %d %s", r.Method, r.URL.String(), w2.StatusCode, http.StatusText(w2.StatusCode))
	})
}

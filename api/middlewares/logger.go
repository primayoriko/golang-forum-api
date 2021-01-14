package middlewares

import (
	"net/http"

	"github.com/gorilla/context"
	"gitlab.com/hydra/forum-api/api/logger"
)

// Log is a middleware method for writing log for every passed request
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Format: type[info/error/fatal/panic] username Host:path?query statuscode message(opt)
		username := context.Get(r, "username").(string)
		id := context.Get(r, "id").(uint32)

		w2 := logger.NewLoggingResponseWriter(w)
		next.ServeHTTP(w2, r)

		log := logger.GetInstance()
		// log.Infof("%s %s, response %d %s", r.Method, r.URL.String(),
		// 	w2.StatusCode, http.StatusText(w2.StatusCode))
		log.WriteLog(id, username, r.Method, r.URL.String(),
			w2.StatusCode, http.StatusText(w2.StatusCode))
	})
}

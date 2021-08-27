package middlewares

import (
	"log"
	"net/http"
)

// EntryLog is a logging middleware stand in front of all Handlers and HandleFunc
type EntryLog struct{}

// Apply logs request before passing it to http.Handler
func (e *EntryLog) Apply(handler http.Handler) http.HandlerFunc {
	return e.ApplyFn(handler.ServeHTTP)
}

// ApplyFn logs request before passing it to http.HandleFunc
func (e *EntryLog) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := "%s %s | %+v\n"
		log.Printf(format, r.Method, r.URL.Path, r.URL.Query().Encode())
		next(w, r)
	}
}

package middlewares

import (
	"log"
	"net/http"
)

type EntryLog struct{}

func (e *EntryLog) Apply(handler http.Handler) http.HandlerFunc {
	return e.ApplyFn(handler.ServeHTTP)
}

func (e *EntryLog) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := "%s %s | %+v\n"
		log.Printf(format, r.Method, r.URL.Path, r.URL.Query().Encode())
		next(w, r)
	}
}

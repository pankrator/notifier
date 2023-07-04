package http

import (
	"log"
	"net/http"
)

func recoverHandler(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if perr := recover(); perr != nil {
				log.Printf("recover from panic: %s", perr)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		f(w, r)
	}
}

package api

import (
	"net/http"
)

// RecoverMiddleware offers a middleware that catches panics.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(writer, req)
	})
}

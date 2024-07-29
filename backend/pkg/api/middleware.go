package api

import (
	"net/http"
)

// EnableCORS middleware.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if req.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
		}

		next.ServeHTTP(writer, req)
	})
}

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

package middleware

import (
	"mime"
	"net/http"
)

func EnforceJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType == "" {
			http.Error(w, "Empty Content-Type header. Content-Type header must be application/json", http.StatusBadRequest)
			return
		} else {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}

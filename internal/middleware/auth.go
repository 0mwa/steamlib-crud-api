package middleware

import (
	"Crud-Api/internal/repository"
	"net/http"
)

func AuthHandler(next http.Handler, sessRepo repository.Sessions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("token")
		ok, err := sessRepo.CheckToken(headerToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

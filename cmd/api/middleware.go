package main

import (
	"net/http"
	"strings"
)

func (a *app) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No valid token", http.StatusBadRequest)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Malformed request", http.StatusBadRequest)
			return
		}

		token := parts[1]
		exist, err := a.auth.ValidateUser(r.Context(), token)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exist {
			http.Error(w, "User is not logged in", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

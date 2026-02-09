package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (a *app) authTokenMiddleware(next http.Handler) http.Handler {
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
		user, err := a.auth.AuthenticateToken(r.Context(), token)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could find user with that token", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

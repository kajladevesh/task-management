package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"task_management/task_service/src/internal/adaptors/external"
)

type AuthMiddleware struct {
	userClient *external.UserServiceClient
}

func NewAuthMiddleware(client *external.UserServiceClient) *AuthMiddleware {
	return &AuthMiddleware{userClient: client}
}

func (a *AuthMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// 1. Try Authorization header: "Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
			fmt.Println("Token found in Authorization header:", token)
		}

		if token == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		// 3. If still empty, reject
		if token == "" {
			http.Error(w, "Unauthorized - no token provided", http.StatusUnauthorized)
			return
		}

		// 4. Validate token with user service (gRPC)
		userID, valid := a.userClient.ValidateToken(token)
		// fmt.Println("Token valid:", valid, "UserID:", userID)

		if !valid {
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		// Add userID to request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

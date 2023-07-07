package grpc_server

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
)

func AuthMiddleware(handler http.Handler, a interfaces.IAuth, s *db.Store) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get the Authorization header value
		authHeader := req.Header.Get("Authorization")

		// Trim leading and trailing whitespace
		authHeader = strings.TrimSpace(authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(res, "missing bearer token", http.StatusUnauthorized)
			return
		}

		// Extract the token by splitting the header value
		parts := strings.Split(authHeader, " ")

		if len(parts) < 2 {
			http.Error(res, "missing bearer token", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		payload, err := a.VerifyToken(token)

		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}

		// Load user in context
		user, err := s.GetUserById(req.Context(), payload.UserID)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(res, "invalid token", http.StatusUnauthorized)
				return
			}

			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), "user", &user))

		// Call the next handler
		handler.ServeHTTP(res, req)
	})

}

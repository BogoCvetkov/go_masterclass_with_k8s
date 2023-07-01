package grpc_server

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/BogoCvetkov/go_mastercalss/db"
	models "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"google.golang.org/grpc/metadata"
)

func AuthenticateUser(c context.Context, a interfaces.IAuth, s *db.Store) (*models.User, error) {
	var authHeader string

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, fmt.Errorf("missing authorization header")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader = values[0]

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("missing bearer token")
	}

	// Extract the token by splitting the header value
	parts := strings.Split(authHeader, " ")

	if len(parts) < 2 {
		return nil, fmt.Errorf("missing bearer token")
	}

	token := parts[1]

	payload, err := a.VerifyToken(token)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	// Load user in context
	user, err := s.GetUserById(c, payload.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid token")
		}

		return nil, fmt.Errorf("authorization failed")
	}

	return &user, nil
}

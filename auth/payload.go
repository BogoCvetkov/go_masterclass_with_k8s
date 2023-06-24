package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	TokenID uuid.UUID `json:"token_id"`
	UserID  int64     `json:"user_id"`
	jwt.RegisteredClaims
}

// Valid checks if the token payload is valid or not
func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiresAt.Time) {
		return errors.New("token has expired")
	}
	return nil
}

func NewTokenPayload(user_id int64, duration time.Duration) (*TokenPayload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	p := &TokenPayload{
		id,
		user_id,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "master_class",
		},
	}

	return p, nil
}

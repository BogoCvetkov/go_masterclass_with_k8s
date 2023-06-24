package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secretKey string
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{
		secretKey: secret,
	}
}

func (a *JWTAuth) GenerateToken(claims *TokenPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JWTAuth) VerifyToken(t string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(t, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		// validate the alg that you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(a.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

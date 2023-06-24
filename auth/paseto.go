package auth

import (
	"github.com/o1egl/paseto"
)

type PasetoAuth struct {
	secretKey string
}

func NewPasetoAuth(secret string) *PasetoAuth {
	return &PasetoAuth{
		secretKey: secret,
	}
}

func (a *PasetoAuth) GenerateToken(claims *TokenPayload) (string, error) {
	token, err := paseto.NewV2().Encrypt([]byte(a.secretKey), claims, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *PasetoAuth) VerifyToken(t string) (*TokenPayload, error) {
	payload := &TokenPayload{}

	err := paseto.NewV2().Decrypt(t, []byte(a.secretKey), payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

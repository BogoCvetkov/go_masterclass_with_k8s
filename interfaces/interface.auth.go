package interfaces

import "github.com/BogoCvetkov/go_mastercalss/auth"

type IAuth interface {
	GenerateToken(claims *auth.TokenPayload) (string, error)
	VerifyToken(t string) (*auth.TokenPayload, error)
}

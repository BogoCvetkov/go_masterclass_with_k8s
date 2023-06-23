package db

import (
	"fmt"

	"github.com/lib/pq"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

func ErrorCode(err error) string {
	// var pgErr *pgconn.PgError
	// if errors.As(err, &pgErr) {
	// 	return pgErr.Code
	// }
	// return ""

	if pqErr, ok := err.(*pq.Error); ok {
		fmt.Println(pqErr)
		return string(pqErr.Code)
	}

	return ""
}

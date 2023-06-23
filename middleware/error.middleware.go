package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Description string `json:"description,omitempty"`
	StatusCode  int    `json:"status_code"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("description: %s", e.Description)
}

func HandleErr(c *gin.Context, msg string, code int) {
	err := ApiError{Description: msg, StatusCode: code}
	c.Error(err)
}

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors[0]
		switch e := err.Err.(type) {
		case ApiError:
			c.JSON(e.StatusCode, e)
			return
		default:
			msg := ApiError{Description: "Unexpected error occured", StatusCode: http.StatusInternalServerError}
			c.JSON(msg.StatusCode, msg)
			return
		}
	}
}

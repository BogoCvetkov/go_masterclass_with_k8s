package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/BogoCvetkov/go_mastercalss/db"
	models "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(a interfaces.IAuth, s db.IStore) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Get the Authorization header value
		authHeader := c.GetHeader("Authorization")

		// Trim leading and trailing whitespace
		authHeader = strings.TrimSpace(authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "missing bearer token"})
			return
		}

		// Extract the token by splitting the header value
		parts := strings.Split(authHeader, " ")

		if len(parts) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "missing bearer token"})
			return
		}

		token := parts[1]

		payload, err := a.VerifyToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}

		// Load user in context
		user, err := s.GetUserById(c, payload.UserID)

		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "authorization failed"})
			return
		}

		c.Set("user", user)

		c.Next()
	}

}

func GetReqUser(c *gin.Context) *models.User {

	user, _ := c.Get("user")
	result, _ := user.(models.User)

	return &result
}

package api

import (
	"database/sql"
	"net/http"
	"time"

	controller "github.com/BogoCvetkov/go_mastercalss/api/controller/types"
	"github.com/BogoCvetkov/go_mastercalss/auth"

	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	m "github.com/BogoCvetkov/go_mastercalss/middleware"
	"github.com/BogoCvetkov/go_mastercalss/util"
	"github.com/gin-gonic/gin"
)

type userController struct {
	*apiController
}

func (ctr *userController) CreateUser(c *gin.Context) {
	var data controller.CreateUserParams

	// Validate data
	if err := c.ShouldBindJSON(&data); err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Encrypt password
	hash, err := util.HashPassword(data.Password)
	if err != nil {
		m.HandleErr(c, "Failed to hash password", http.StatusBadRequest)
		return
	}

	document := db.CreateUserParams{
		Username:       data.Username,
		FullName:       data.FullName,
		Email:          data.Email,
		HashedPassword: hash,
	}

	// Create new user
	user, err := ctr.server.GetStore().CreateUserTrx(c, document, ctr.server.GetAsync())
	if err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	result := filterUser(user)

	c.JSON(http.StatusOK, result)

}

func (ctr *userController) LoginUser(c *gin.Context) {
	var data controller.LoginUserParams

	// Validate data
	if err := c.ShouldBindJSON(&data); err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Find user
	user, err := ctr.server.GetStore().GetUser(c, data.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			m.HandleErr(c, "User not found", http.StatusNotFound)
			return
		}

		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify password
	err = util.CheckPassword(data.Password, user.HashedPassword)
	if err != nil {
		m.HandleErr(c, "Invalid credentials", http.StatusBadRequest)
		return
	}

	// Prepare access token payload
	p1, err := auth.NewTokenPayload(user.ID, ctr.server.GetConfig().TokenDuration)
	if err != nil {
		m.HandleErr(c, "Failed generating token", http.StatusBadRequest)
		return
	}

	// Generate access token
	token, err := ctr.server.GetAuth().GenerateToken(p1)
	if err != nil {
		m.HandleErr(c, "Failed generating token", http.StatusBadRequest)
		return
	}

	// Prepare REFRESH token payload
	p2, err := auth.NewTokenPayload(user.ID, ctr.server.GetConfig().RTokenDuration)
	if err != nil {
		m.HandleErr(c, "Failed generating refresh token", http.StatusBadRequest)
		return
	}

	// Generate REFRESH token
	rtoken, err := ctr.server.GetAuth().GenerateToken(p2)
	if err != nil {
		m.HandleErr(c, "Failed generating refresh token", http.StatusBadRequest)
		return
	}

	// Store session
	arg := db.CreateSessionParams{
		ID:           p2.TokenID,
		UserID:       user.ID,
		RefreshToken: rtoken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    p2.ExpiresAt.Time,
	}
	_, err = ctr.server.GetStore().CreateSession(c, arg)
	if err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	filtered := filterUser(&user)

	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": rtoken, "user": filtered})
}

func (ctr *userController) RefreshToken(c *gin.Context) {
	var data controller.RefreshTokenParams

	// Validate data
	if err := c.ShouldBindJSON(&data); err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify refresh token
	payload, err := ctr.server.GetAuth().VerifyToken(data.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		return
	}

	// Check user session data
	session, err := ctr.server.GetStore().GetSession(c, payload.TokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "authorization failed"})
		return
	}

	if payload.UserID != session.UserID || data.Token != session.RefreshToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
		return
	}

	if session.IsBlocked {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "expired token"})
		return
	}

	// Prepare access token payload
	p, err := auth.NewTokenPayload(payload.UserID, ctr.server.GetConfig().TokenDuration)
	if err != nil {
		m.HandleErr(c, "Failed generating token", http.StatusBadRequest)
		return
	}

	// Generate access token
	token, err := ctr.server.GetAuth().GenerateToken(p)
	if err != nil {
		m.HandleErr(c, "Failed generating token", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func filterUser(u *db.User) *controller.CreateUserReponse {

	result := &controller.CreateUserReponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
	}

	return result

}

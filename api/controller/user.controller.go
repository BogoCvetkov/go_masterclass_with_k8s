package api

import (
	"database/sql"
	"net/http"

	controller "github.com/BogoCvetkov/go_mastercalss/api/controller/types"
	"github.com/BogoCvetkov/go_mastercalss/auth"

	db_util "github.com/BogoCvetkov/go_mastercalss/db"
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
	user, err := ctr.server.GetStore().CreateUser(c, document)
	if err != nil {

		if db_util.ErrorCode(err) == db_util.UniqueViolation {
			m.HandleErr(c, "Email already exists", http.StatusBadRequest)
			return
		}

		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	result := filterUser(&user)

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
	p, err := auth.NewTokenPayload(user.ID, ctr.server.GetConfig().TokenDuration)
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

	filtered := filterUser(&user)

	c.JSON(http.StatusOK, gin.H{"token": token, "user": filtered})
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
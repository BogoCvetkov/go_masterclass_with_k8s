package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	controller "github.com/BogoCvetkov/go_mastercalss/controller/types"
	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/gin-gonic/gin"
)

type accController struct {
	*apiController
}

func (ctr *accController) CreateAccount(c *gin.Context) {
	var data controller.CreateAccountParams

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document := db.CreateAccountParams{
		Owner:    data.Owner,
		Currency: data.Currency,
		Balance:  0,
	}

	acc, err := ctr.store.CreateAccount(c, document)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (ctr *accController) GetAccount(c *gin.Context) {

	param := c.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := ctr.store.GetAccount(c, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (ctr *accController) ListAccounts(c *gin.Context) {

	var q controller.ListAccountQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := db.ListAccountsParams{
		Limit:  q.Limit,
		Offset: (q.Page - 1) * q.Limit,
	}

	acc, err := ctr.store.ListAccounts(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}

package api

import (
	"database/sql"
	"net/http"
	"strconv"

	controller "github.com/BogoCvetkov/go_mastercalss/api/controller/types"
	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	m "github.com/BogoCvetkov/go_mastercalss/middleware"
	"github.com/gin-gonic/gin"
)

type accController struct {
	*apiController
}

func (ctr *accController) CreateAccount(c *gin.Context) {
	var data controller.CreateAccountParams

	if err := c.ShouldBindJSON(&data); err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	reqU := m.GetReqUser(c)

	document := db.CreateAccountParams{
		Owner:    reqU.ID,
		Currency: data.Currency,
		Balance:  0,
	}

	acc, err := ctr.server.GetStore().CreateAccount(c, document)
	if err != nil {
		m.HandleErr(c, "Failed to create account", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (ctr *accController) GetAccount(c *gin.Context) {

	param := c.Param("id")

	reqU := m.GetReqUser(c)

	id, err := strconv.Atoi(param)
	if err != nil {
		m.HandleErr(c, "Id not an integer", http.StatusBadRequest)
		return
	}

	params := db.GetAccountByOwnerParams{
		ID:    int64(id),
		Owner: reqU.ID,
	}

	acc, err := ctr.server.GetStore().GetAccountByOwner(c, params)
	if err != nil {
		if err == sql.ErrNoRows {
			m.HandleErr(c, "Account not found", http.StatusNotFound)
			return
		}

		m.HandleErr(c, "Failed to get account", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (ctr *accController) ListAccounts(c *gin.Context) {

	var q controller.ListAccountQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	reqU := m.GetReqUser(c)

	query := db.ListAccountsParams{
		Owner:  reqU.ID,
		Limit:  q.Limit,
		Offset: (q.Page - 1) * q.Limit,
	}

	acc, err := ctr.server.GetStore().ListAccounts(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			m.HandleErr(c, "No accounts found", http.StatusNotFound)
			return
		}

		m.HandleErr(c, "Failed to get accounts", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, acc)
}

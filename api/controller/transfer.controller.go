package api

import (
	"database/sql"
	"fmt"
	"net/http"

	controller "github.com/BogoCvetkov/go_mastercalss/api/controller/types"
	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	m "github.com/BogoCvetkov/go_mastercalss/middleware"
	"github.com/gin-gonic/gin"
)

type transferController struct {
	*apiController
}

func (ctr *transferController) CreateTransfer(c *gin.Context) {
	var data controller.CreateTransferParams

	if err := c.ShouldBindJSON(&data); err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	reqU := m.GetReqUser(c)

	acc, ok := ctr.validAccount(c, data.FromAccountID, data.Currency)
	if !ok {
		return
	}

	if acc.Owner != reqU.ID {
		m.HandleErr(c, "from acc does not belong to caller", http.StatusNotFound)
		return
	}

	if _, ok := ctr.validAccount(c, data.ToAccountID, data.Currency); !ok {
		return
	}

	document := db.CreateTransferParams{
		FromAccountID: data.FromAccountID,
		ToAccountID:   data.ToAccountID,
		Amount:        data.Amount,
	}

	result, err := ctr.server.GetStore().TransferTrx(c, document)
	if err != nil {
		m.HandleErr(c, err.Error(), http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctr *transferController) validAccount(c *gin.Context, accountID int64, currency string) (*db.Account, bool) {
	account, err := ctr.server.GetStore().GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			m.HandleErr(c, "Account not found", http.StatusNotFound)
			return nil, false
		}

		m.HandleErr(c, "Failed to get account", http.StatusBadRequest)
		return nil, false
	}

	if account.Currency != currency {
		msg := fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		m.HandleErr(c, msg, http.StatusBadRequest)
		return nil, false
	}

	return &account, true
}

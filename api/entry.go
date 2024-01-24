package api

import (
	"net/http"

	db "github.com/Crazypointer/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required,min=1"`
}

func (server Server) createEntry(ctx *gin.Context) {
	var req CreateEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount}
	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

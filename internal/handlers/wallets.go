package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"test-task/internal/dto"
	"test-task/internal/services"
)

type walletService interface {
	GetBalance(ctx context.Context, id string) (float64, error)
	RunOperation(ctx context.Context, operation services.WalletOperation) error
}

type WalletHandler struct {
	service walletService
}

func NewWalletHandler(service walletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) GetBalance(ctx *gin.Context) {

	walletID := ctx.Param("id")

	if walletID == "" {
		ctx.JSON(400, dto.ErrorResponse{Error: "wallet ID is required"})
		return
	}

	if _, err := uuid.Parse(walletID); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: "wallet ID must be uuid"})
		return
	}

	balance, err := h.service.GetBalance(ctx, walletID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(200, dto.WalletBalance{Balance: balance})
}

func (h *WalletHandler) RunOperation(ctx *gin.Context) {

	var dtoOp dto.WalletOperation
	if err := ctx.ShouldBindJSON(&dtoOp); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	op, err := services.NewWalletOperation(dtoOp.WalledID, dtoOp.OperationType, dtoOp.Amount)
	if err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = h.service.RunOperation(ctx, *op)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
}

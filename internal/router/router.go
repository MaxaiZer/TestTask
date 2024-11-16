package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"test-task/internal/dto"
	errs "test-task/internal/errors"
	"test-task/internal/handlers"
)

func Setup(engine *gin.Engine, walletHandler *handlers.WalletHandler) {

	engine.Use(gin.Recovery())
	engine.Use(errorHandler)

	engine.GET("/api/v1/wallets/:id", walletHandler.GetBalance)
	engine.POST("/api/v1/wallet", walletHandler.RunOperation)
}

func errorHandler(ctx *gin.Context) {

	ctx.Next()

	if len(ctx.Errors) > 0 {
		err := ctx.Errors.Last()

		if errors.Is(err, errs.NotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, errs.UnsupportedOperation) {
			ctx.AbortWithStatusJSON(http.StatusNotImplemented, dto.ErrorResponse{Error: err.Error()})
			logError(ctx, err)
		} else if errors.Is(err, errs.InsufficientBalance) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, errs.TooManyRequests) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: err.Error()})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
			logError(ctx, err)
		}
	}
}

func logError(ctx *gin.Context, err error) {
	log.Errorf("%s %s error: %s", ctx.Request.Method, ctx.Request.URL.Path, err)
}

package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
	"test-task/internal/config"
	"test-task/internal/handlers"
	"test-task/internal/repositories"
	"test-task/internal/router"
	"test-task/internal/services"
)

func main() {

	cfg := config.Get()

	dbContext, err := repositories.NewDbContext(cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("error create dbContext: %v", err)
		return
	}

	err = dbContext.Migrate()
	if err != nil {
		log.Fatalf("error run database migration: %v", err)
		return
	}
	log.Infof("database migration complete")

	walletRepository := repositories.NewWalletsRepository(dbContext.DB)
	walletService := services.NewWalletsService(walletRepository)
	defer walletService.Close()
	walletHandler := handlers.NewWalletHandler(walletService)

	gin.SetMode(cfg.Mode)
	ginEngine := gin.New()

	router.Setup(ginEngine, walletHandler)

	log.Errorf(ginEngine.Run(":" + strconv.Itoa(cfg.Port)).Error())
}

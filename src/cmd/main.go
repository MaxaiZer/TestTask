package main

import (
	"project/src/config"
	"project/src/server"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Get()

	gin.SetMode(cfg.Mode)
	ginEngine := gin.Default()

	server.InitializeRoutes(ginEngine)

	ginEngine.Run(":" + strconv.Itoa(cfg.Port))
}

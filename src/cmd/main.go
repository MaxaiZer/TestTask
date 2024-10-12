// @title TestTask API
// @version 1.0
// @description Already created users with id 1 and 2
// @host localhost:8080
// @basePath /
// @schemes http

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

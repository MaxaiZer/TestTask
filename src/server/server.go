package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"project/docs"
	"project/src/config"
	"project/src/errors"
	"project/src/handlers"
	"strconv"
)

func InitializeRoutes(engine *gin.Engine) {
	engine.Use(errorHandler)
	handlers.Register(engine)

	if config.Get().Mode == "debug" {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Host = "localhost:" + strconv.Itoa(config.Get().Port)
	}
}

func errorHandler(c *gin.Context) {

	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last()

		switch e := err.Err.(type) {
		case errors.PublicError:
			c.AbortWithStatusJSON(e.Code, e)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}
	}
}

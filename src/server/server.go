package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/src/errors"
	"project/src/handlers"
)

func InitializeRoutes(engine *gin.Engine) {
	engine.Use(errorHandler)
	handlers.Register(engine)
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

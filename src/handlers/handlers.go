package handlers

import (
	"net/http"
	"project/src/config"
	"project/src/dto"
	"project/src/errors"
	"project/src/repositories"
	"project/src/services"

	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	e.POST("/auth/tokens", createTokens)
	e.POST("/auth/refresh", refreshTokens)
}

// @Summary Create new tokens
// @Description Creates new tokens for a user based on the user_id query parameter
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} dto.TokenPair
// @Failure 400 {object} errors.PublicError
// @Router /auth/tokens [post]
func createTokens(c *gin.Context) {

	userID := c.Query("user_id")
	if userID == "" {
		_ = c.Error(errors.PublicError{
			Code:    http.StatusBadRequest,
			Message: "Missing user_id query parameter",
		})
		return
	}

	service, err := createAuthService()
	if err != nil {
		_ = c.Error(err)
		return
	}

	tokens, err := service.CreateTokens(c.Request.Context(), c.ClientIP(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary Refresh tokens
// @Description Refreshes tokens based on the provided token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param pair body dto.TokenPair true "Access and refresh tokens"
// @Success 200 {object} dto.TokenPair
// @Failure 400 {object} errors.PublicError
// @Failure 401 {object} errors.PublicError
// @Router /auth/refresh [post]
func refreshTokens(c *gin.Context) {

	var pair dto.TokenPair

	if err := c.ShouldBindJSON(&pair); err != nil {
		_ = c.Error(errors.PublicError{Code: http.StatusBadRequest, Message: "Invalid json format"})
		return
	}

	service, err := createAuthService()
	if err != nil {
		_ = c.Error(err)
		return
	}

	tokens, err := service.RefreshTokens(c.Request.Context(), c.ClientIP(), pair)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func createAuthService() (*services.AuthService, error) {

	cfg := config.Get()

	jwt, err := services.NewJwtService(&cfg.JWT)
	if err != nil {
		return nil, err
	}

	repository, err := repositories.NewPostgresUserRepository(cfg.DB.ConnectionString)
	if err != nil {
		return nil, err
	}

	service := services.NewAuthService(jwt, repository, &services.NotifyService{})
	return service, nil
}

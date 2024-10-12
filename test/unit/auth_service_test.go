package unit

import (
	"context"
	"project/src/config"
	"project/src/dto"
	"project/src/entities"
	"project/src/services"
	"testing"
	"time"
)

func Test_CreateTokens_WhenUserDoesntExist_ShouldReturnError(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	user := entities.User{ID: "0", Email: "mail@mail.com", RefreshToken: "", RefreshTokenExpiryTime: time.Time{}}

	mockRepository := NewMockUserRepository([]entities.User{user})

	authService := services.NewAuthService(jwtService, mockRepository, &MockNotifyService{})

	_, err = authService.CreateTokens(context.Background(), "127.0.0.1", "1")
	if err == nil {
		t.Fatalf("createTokens success when expecting error")
	}
}

func Test_CreateTokens_WhenUserExist_ShouldSuccess(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	user := entities.User{ID: "0", Email: "mail@mail.com", RefreshToken: "", RefreshTokenExpiryTime: time.Time{}}

	mockRepository := NewMockUserRepository([]entities.User{user})

	authService := services.NewAuthService(jwtService, mockRepository, &MockNotifyService{})

	_, err = authService.CreateTokens(context.Background(), "127.0.0.1", "0")
	if err != nil {
		t.Fatalf("failed to create tokens: %v", err)
	}
}

func Test_RefreshTokens_WhenValidData_ShouldSuccess(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	user := entities.User{ID: "0", Email: "mail@mail.com", RefreshToken: "0", RefreshTokenExpiryTime: time.Now()}
	userIp := "127.0.0.1"

	mockRepository := NewMockUserRepository([]entities.User{user})

	authService := services.NewAuthService(jwtService, mockRepository, &MockNotifyService{})

	userTokens, err := authService.CreateTokens(context.Background(), userIp, "0")
	if err != nil {
		t.Fatalf("failed to create tokens: %v", err)
	}

	_, err = authService.RefreshTokens(context.Background(), userIp, *userTokens)
	if err != nil {
		t.Fatalf("failed to refresh tokens: %v", err)
	}
}

func Test_RefreshTokens_WhenIpChanged_ShouldSuccessAndCallNotifier(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	userEmail := "mail@mail.com"
	user := entities.User{ID: "0", Email: userEmail, RefreshToken: "0", RefreshTokenExpiryTime: time.Now()}
	userIp := "127.0.0.1"
	newIp := "8.8.8.8"

	mockRepository := NewMockUserRepository([]entities.User{user})
	mockNotifyService := MockNotifyService{}

	authService := services.NewAuthService(jwtService, mockRepository, &mockNotifyService)

	userTokens, err := authService.CreateTokens(context.Background(), userIp, "0")
	if err != nil {
		t.Fatalf("failed to create tokens: %v", err)
	}

	_, err = authService.RefreshTokens(context.Background(), newIp, *userTokens)
	if err != nil {
		t.Fatalf("failed to refresh tokens: %v", err)
	}

	if !mockNotifyService.IsNotified(userEmail) {
		t.Fatalf("user not notified about ip change")
	}
}

func Test_RefreshTokens_WhenInvalidAccessToken_ShouldReturnError(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	user := entities.User{ID: "0", Email: "mail@mail.com", RefreshToken: "0", RefreshTokenExpiryTime: time.Now()}
	userIp := "127.0.0.1"

	mockRepository := NewMockUserRepository([]entities.User{user})

	authService := services.NewAuthService(jwtService, mockRepository, &MockNotifyService{})

	userTokens, err := authService.CreateTokens(context.Background(), userIp, "0")
	if err != nil {
		t.Fatalf("failed to create tokens: %v", err)
	}

	userTokens.AccessToken += "s"

	_, err = authService.RefreshTokens(context.Background(), userIp, *userTokens)
	if err == nil {
		t.Fatalf("refreshTokens success when expecting error")
	}
}

func Test_RefreshTokens_WhenInvalidRefreshToken_ShouldReturnError(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	user := entities.User{ID: "0", Email: "mail@mail.com", RefreshToken: "0", RefreshTokenExpiryTime: time.Now()}
	userIp := "127.0.0.1"

	mockRepository := NewMockUserRepository([]entities.User{user})

	authService := services.NewAuthService(jwtService, mockRepository, &MockNotifyService{})

	userTokens, err := authService.CreateTokens(context.Background(), userIp, "0")
	if err != nil {
		t.Fatalf("failed to create tokens: %v", err)
	}

	userTokens.RefreshToken += "s"

	_, err = authService.RefreshTokens(context.Background(), userIp, *userTokens)
	if err == nil {
		t.Fatalf("refreshTokens success when expecting error")
	}
}

func Test_RefreshTokens_WhenAlreadyUsedRefreshToken_ShouldReturnError(t *testing.T) {

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	user := entities.User{ID: "0", Email: "mail@mail.com", RefreshToken: "0", RefreshTokenExpiryTime: time.Now()}
	userIp := "127.0.0.1"

	mockRepository := NewMockUserRepository([]entities.User{user})

	authService := services.NewAuthService(jwtService, mockRepository, &MockNotifyService{})

	userTokens1, _ := authService.CreateTokens(context.Background(), userIp, "0")

	userTokens2, err := authService.RefreshTokens(context.Background(), userIp, *userTokens1)
	if err != nil {
		t.Fatalf("failed to refresh tokens: %v", err)
	}

	userTokens3 := dto.TokenPair{AccessToken: userTokens2.AccessToken, RefreshToken: userTokens1.RefreshToken}

	_, err = authService.RefreshTokens(context.Background(), userIp, userTokens3)
	if err == nil {
		t.Fatalf("refreshTokens success when expecting error")
	}
}

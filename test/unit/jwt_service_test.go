package unit

import (
	"log"
	"project/src/config"
	"project/src/services"
	"testing"
)

func Test_CreateTokenPair_ShouldSuccessWithValidAccessToken(t *testing.T) {

	claims := services.ExtraClaims{UserId: "12345", Ip: "127.0.0.1"}

	jwtService, err := services.NewJwtService(&config.Get().JWT)
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	pair, err := jwtService.CreateTokenPair(claims)
	if err != nil {
		t.Fatalf("failed to create token pair: %v", err)
	}

	if pair.AccessToken == "" {
		t.Fatalf("access token is empty")
	}

	if pair.RefreshToken == "" {
		t.Fatalf("refresh token is empty")
	}

	claimsFromToken, err := jwtService.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	if claims != *claimsFromToken {
		t.Fatalf("resulting after validation claims are different from original: %v", claimsFromToken)
	}
}

func Test_ValidateAccessToken_WhenInvalidIssuer_ShouldReturnError(t *testing.T) {

	jwtService1, err := services.NewJwtService(&config.Get().JWT)

	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	jwtService2 := CreateJwtServiceWithIssuerAndAudience("AnotherIssuer", config.Get().JWT.Audience)

	pair, err := jwtService2.CreateTokenPair(services.ExtraClaims{UserId: "12345", Ip: "127.0.0.1"})
	if err != nil {
		t.Fatalf("failed to create token pair: %v", err)
	}

	_, err = jwtService1.ValidateAccessToken(pair.AccessToken)
	if err == nil {
		t.Fatalf("ValidateAccessToken didnt return error")
	}
}

func Test_ValidateAccessToken_WhenInvalidAudience_ShouldReturnError(t *testing.T) {

	jwtService1, err := services.NewJwtService(&config.Get().JWT)

	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	jwtService2 := CreateJwtServiceWithIssuerAndAudience(config.Get().JWT.Issuer, "AnotherAudience")

	pair, err := jwtService2.CreateTokenPair(services.ExtraClaims{UserId: "12345", Ip: "127.0.0.1"})
	if err != nil {
		t.Fatalf("failed to create token pair: %v", err)
	}

	_, err = jwtService1.ValidateAccessToken(pair.AccessToken)
	if err == nil {
		t.Fatalf("ValidateAccessToken didnt return error")
	}
}

func CreateJwtServiceWithIssuerAndAudience(issuer string, audience string) *services.JwtService {

	jwt, err := services.NewJwtService(&config.JWTConfig{
		AccessLifetime:  3600,
		RefreshLifetime: 7200,
		SecretKey:       config.Get().JWT.SecretKey,
		Issuer:          issuer,
		Audience:        audience,
	})

	if err != nil {
		log.Fatalf("failed to create jwt service: %v", err)
	}

	return jwt
}

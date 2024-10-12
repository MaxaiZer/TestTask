package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"project/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretKey          string
	issuer             string
	audience           string
	accessLifetime     int
	refreshLifetime    int
	signingMethod      jwt.SigningMethod
	refreshTokenLength int
}

type ExtraClaims struct {
	UserId string `json:"user_id"`
	Ip     string `json:"ip"`
}

type Tokens struct {
	AccessToken            string
	RefreshToken           string
	RefreshTokenExpiryTime time.Time
}

type fullClaims struct {
	jwt.RegisteredClaims
	ExtraClaims
}

func NewJwtService(cfg *config.JWTConfig) (*JwtService, error) {

	if cfg.AccessLifetime <= 0 {
		return nil, fmt.Errorf("invalid Jwt AccessLifetime: %d", cfg.AccessLifetime)
	}

	if cfg.RefreshLifetime <= 0 {
		return nil, fmt.Errorf("invalid Jwt RefreshLifetime: %d", cfg.RefreshLifetime)
	}

	return &JwtService{
		secretKey:          cfg.SecretKey,
		signingMethod:      jwt.SigningMethodHS512,
		accessLifetime:     cfg.AccessLifetime,
		refreshLifetime:    cfg.RefreshLifetime,
		issuer:             cfg.Issuer,
		audience:           cfg.Audience,
		refreshTokenLength: 32,
	}, nil
}

func (s JwtService) CreateTokenPair(accessTokenExtraClaims ExtraClaims) (*Tokens, error) {
	access, err := s.createAccessToken(accessTokenExtraClaims)
	if err != nil {
		return nil, err
	}

	refresh, err := s.createRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshExpiryTime := time.Now().Add(time.Duration(s.refreshLifetime) * time.Second)

	return &Tokens{AccessToken: access, RefreshToken: refresh, RefreshTokenExpiryTime: refreshExpiryTime}, nil
}

func (s JwtService) ValidateAccessToken(tokenString string) (*ExtraClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &fullClaims{}, func(token *jwt.Token) (any, error) {

		if token.Method != s.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*fullClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if s.issuer != claims.Issuer || len(claims.Audience) == 0 || s.audience != claims.Audience[0] {
		return nil, fmt.Errorf("invalid audience or issuer")
	}

	return &claims.ExtraClaims, nil
}

func (s JwtService) createAccessToken(extraClaims ExtraClaims) (string, error) {

	claims := fullClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.accessLifetime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ExtraClaims: extraClaims,
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)

	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s JwtService) createRefreshToken() (string, error) {

	bytes := make([]byte, s.refreshTokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

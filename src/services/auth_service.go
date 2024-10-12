package services

import (
	"context"
	"net/http"
	"project/src/dto"
	"project/src/errors"
	"project/src/interfaces"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository interfaces.UserRepository
	notifyService  interfaces.NotifyService
	jwtService     *JwtService
}

func NewAuthService(
	jwtService *JwtService,
	userRepository interfaces.UserRepository,
	notifyService interfaces.NotifyService,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
		notifyService:  notifyService,
	}
}

func (s *AuthService) CreateTokens(ctx context.Context, clientIp string, userId string) (*dto.TokenPair, error) {

	user, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.PublicError{Code: http.StatusBadRequest, Message: "Invalid user id"}
	}

	tokenPair, err := s.jwtService.CreateTokenPair(ExtraClaims{Ip: clientIp, UserId: userId})
	if err != nil {
		return nil, err
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokenPair.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.RefreshToken = string(hashedRefreshToken)
	user.RefreshTokenExpiryTime = tokenPair.RefreshTokenExpiryTime

	err = s.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, clientIp string, tokens dto.TokenPair) (*dto.TokenPair, error) {

	claims, err := s.jwtService.ValidateAccessToken(tokens.AccessToken)
	if err != nil {
		return nil, errors.PublicError{Code: http.StatusBadRequest, Message: "Invalid access token"}
	}

	user, err := s.userRepository.GetById(ctx, claims.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.PublicError{Code: http.StatusBadRequest, Message: "Invalid user id"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(tokens.RefreshToken)); err != nil {
		return nil, errors.PublicError{Code: http.StatusBadRequest, Message: "Invalid refresh token"}
	}

	if user.RefreshTokenExpiryTime.Compare(time.Now()) < 0 {
		return nil, errors.PublicError{Code: http.StatusUnauthorized, Message: "Refresh token expired"}
	}

	tokenPair, err := s.jwtService.CreateTokenPair(ExtraClaims{Ip: clientIp, UserId: claims.UserId})
	if err != nil {
		return nil, err
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokenPair.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.RefreshToken = string(hashedRefreshToken)

	err = s.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if clientIp != claims.Ip {
		s.notifyService.NotifyAboutIpChange(user.Email)
	}

	return &dto.TokenPair{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken}, nil
}

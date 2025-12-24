package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kanta/backend-challenge/internal/core/domain"
	"github.com/kanta/backend-challenge/internal/core/ports"
)

func GenerateTokenPairWithCache(ctx context.Context, userID string, secret string, cache ports.CachePort) (*domain.TokenPair, error) {
	accessToken, err := generateToken(userID, secret, "access", time.Minute*15)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(userID, secret, "refresh", time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	accessKey := fmt.Sprintf("access:%s", userID)
	if err := cache.SetToken(ctx, accessKey, accessToken, time.Minute*15); err != nil {
		return nil, err
	}

	refreshKey := fmt.Sprintf("refresh:%s", userID)
	if err := cache.SetToken(ctx, refreshKey, refreshToken, time.Hour*24*7); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateToken(userID, secret, tokenType string, duration time.Duration) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshAccessTokenWithCache(ctx context.Context, refreshToken, secret string, cache ports.CachePort) (string, error) {
	claims, err := ParseToken(refreshToken, secret)
	if err != nil {
		return "", err
	}

	if claims.Type != "refresh" {
		return "", errors.New("invalid token type")
	}

	refreshKey := fmt.Sprintf("refresh:%s", claims.UserID)
	storedToken, err := cache.GetToken(ctx, refreshKey)
	if err != nil {
		return "", errors.New("refresh token expired or not found")
	}

	if storedToken != refreshToken {
		return "", errors.New("invalid refresh token")
	}

	accessToken, err := generateToken(claims.UserID, secret, "access", time.Minute*15)
	if err != nil {
		return "", err
	}

	accessKey := fmt.Sprintf("access:%s", claims.UserID)
	if err := cache.SetToken(ctx, accessKey, accessToken, time.Minute*15); err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateAccessToken(tokenStr, secret string) (string, error) {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return "", err
	}

	if claims.Type != "access" {
		return "", errors.New("invalid token type")
	}

	return claims.UserID, nil
}
func ValidateAccessTokenWithCache(ctx context.Context, tokenStr, secret string, cache ports.CachePort) (string, error) {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return "", err
	}

	if claims.Type != "access" {
		return "", errors.New("invalid token type")
	}

	accessKey := fmt.Sprintf("access:%s", claims.UserID)
	storedToken, err := cache.GetToken(ctx, accessKey)
	if err != nil {
		return "", errors.New("access token expired or not found")
	}

	if storedToken != tokenStr {
		return "", errors.New("invalid access token")
	}

	return claims.UserID, nil
}

func ValidateRefreshToken(tokenStr, secret string) (string, error) {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return "", err
	}

	if claims.Type != "refresh" {
		return "", errors.New("invalid token type")
	}

	return claims.UserID, nil
}

func ValidateRefreshTokenWithCache(ctx context.Context, tokenStr, secret string, cache ports.CachePort) (string, error) {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return "", err
	}

	if claims.Type != "refresh" {
		return "", errors.New("invalid token type")
	}

	refreshKey := fmt.Sprintf("refresh:%s", claims.UserID)
	storedToken, err := cache.GetToken(ctx, refreshKey)
	if err != nil {
		return "", errors.New("refresh token expired or not found")
	}

	if storedToken != tokenStr {
		return "", errors.New("invalid refresh token")
	}

	return claims.UserID, nil
}

func RevokeToken(ctx context.Context, userID string, cache ports.CachePort) error {
	accessKey := fmt.Sprintf("access:%s", userID)
	refreshKey := fmt.Sprintf("refresh:%s", userID)

	if err := cache.DeleteToken(ctx, accessKey); err != nil {
		return err
	}

	if err := cache.DeleteToken(ctx, refreshKey); err != nil {
		return err
	}

	return nil
}

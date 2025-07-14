package utils

import (
	"errors"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}

	if config.AppConfig.AccessTokenSecret == "" {
		return "", errors.New("access token secret is not configured")
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.AppConfig.AppName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(config.AppConfig.AccessTokenSecret)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to sign access token: " + err.Error())
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}

	if config.AppConfig.RefreshTokenSecret == "" {
		return "", errors.New("refresh token secret is not configured")
	}

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    config.AppConfig.AppName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(config.AppConfig.RefreshTokenSecret)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to sign refresh token: " + err.Error())
	}

	return tokenString, nil
}

func DecodeAccessToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, errors.New("failed to parse token: " + err.Error())
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func DecodeRefreshToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.RefreshTokenSecret), nil
	})

	if err != nil {
		return "", errors.New("failed to parse refresh token: " + err.Error())
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("invalid refresh token claims")
}

func SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	domain := config.AppConfig.CookieDomain
	c.SetCookie("refreshToken", refreshToken, 7*24*3600, "/", domain, false, false)
}

func ClearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		"refreshToken",
		"",
		-1,
		"/",
		"",
		false,
		false,
	)
}

func SetAccessTokenCookie(c *gin.Context, accessToken string) {
	domain := config.AppConfig.CookieDomain
	c.SetCookie("accessToken", accessToken, 3600, "/", domain, false, false)
}

func ClearAccessTokenCookie(c *gin.Context) {
	c.SetCookie(
		"accessToken",
		"",
		-1,
		"/",
		"",
		false,
		false,
	)
}

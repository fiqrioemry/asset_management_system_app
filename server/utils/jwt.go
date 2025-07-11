package utils

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var accessTokenSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
var refreshTokenSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

func DecodeAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return accessTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid access token")
}

func DecodeRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return refreshTokenSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", errors.New("invalid refresh token")
}

func SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	domain := os.Getenv("COOKIE_DOMAIN")
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", domain, false, false)
}

func ClearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		false,
	)
}

func SetAccessTokenCookie(c *gin.Context, accessToken string) {
	domain := os.Getenv("COOKIE_DOMAIN")
	c.SetCookie("access_token", accessToken, 3600, "/", domain, false, false)
}

func ClearAccessTokenCookie(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		false,
	)
}

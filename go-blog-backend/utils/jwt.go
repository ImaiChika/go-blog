package utils

import (
	"go-blog-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtKey() []byte {
	return []byte(config.GetJWTSecret())
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey())
}
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey(), nil
	})
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, err
}

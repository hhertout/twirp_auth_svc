package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServiceInterface interface {
	Generate(user string) (string, error)
	Verify(tokenString string) (bool, jwt.RegisteredClaims, error)
}

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (j *JWTService) Generate(user string) (string, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return "", errors.New("env variable JWT_SECRET is not set")
	}

	expiresAt := time.Now().Unix() + 3600*24*20

	claims := jwt.RegisteredClaims{
		Issuer:    user,
		ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWTService) Verify(tokenString string) (bool, jwt.RegisteredClaims, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return false, jwt.RegisteredClaims{}, errors.New("env variable JWT_SECRET is not set")
	}

	var claims jwt.RegisteredClaims

	token, err := jwt.NewParser().ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return false, jwt.RegisteredClaims{}, err
	}

	valid := token.Valid

	return valid, claims, nil
}

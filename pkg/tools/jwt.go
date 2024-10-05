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

// NewJWTService creates a new instance of JWTService.
// Returns a pointer to the newly created JWTService.
func NewJWTService() *JWTService {
	return &JWTService{}
}

// Generate creates a JWT token for a given user.
// Uses an environment variable "JWT_SECRET" as the secret key.
// The token expires in 20 days from the time of generation.
// Returns the signed JWT token and an error if any occurs.
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

// Verify checks if a given JWT token is valid.
// Uses an environment variable "JWT_SECRET" as the secret key.
// Returns a boolean indicating if the token is valid, the token claims, and an error if any occurs.
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

package crypto_test

import (
	"os"
	"testing"

	"github.com/hhertout/twirp_auth/lib/crypto"
)

func TestNewJWTService(t *testing.T) {
	jwtService := crypto.NewJWTService()
	if jwtService == nil {
		t.Errorf("expected jwtService to be non-nil")
	}
}

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	jwtService := crypto.NewJWTService()
	token, err := jwtService.Generate("user@example.com")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("expected non-empty token, got empty string")
	}
}

func TestGenerate_NoSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	jwtService := crypto.NewJWTService()
	token, err := jwtService.Generate("user@example.com")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "env variable JWT_SECRET is not set" {
		t.Errorf("expected error 'env variable JWT_SECRET is not set', got %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token, got %v", token)
	}
}

func TestVerifyToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	jwtService := crypto.NewJWTService()
	token, err := jwtService.Generate("user@example.com")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	valid, claims, err := jwtService.Verify(token)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !valid {
		t.Errorf("expected valid token, got invalid")
	}
	if claims.Issuer != "user@example.com" {
		t.Errorf("expected issuer 'user@example.com', got %v", claims.Issuer)
	}
}

func TestVerify_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	jwtService := crypto.NewJWTService()
	invalidToken := "invalid.token.string"

	valid, claims, err := jwtService.Verify(invalidToken)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if valid {
		t.Errorf("expected invalid token, got valid")
	}
	if claims.Issuer != "" {
		t.Errorf("expected empty issuer, got %v", claims.Issuer)
	}
}

func TestVerify_NoSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	jwtService := crypto.NewJWTService()
	token, err := jwtService.Generate("user@example.com")
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	valid, claims, err := jwtService.Verify(token)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "env variable JWT_SECRET is not set" {
		t.Errorf("expected error 'env variable JWT_SECRET is not set', got %v", err)
	}
	if valid {
		t.Errorf("expected invalid token, got valid")
	}
	if claims.Issuer != "" {
		t.Errorf("expected empty issuer, got %v", claims.Issuer)
	}
}

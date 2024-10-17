package crypto_test

import (
	"os"
	"testing"

	"github.com/hhertout/twirp_auth/pkg/crypto"
)

func TestGeneratePassword(t *testing.T) {
	passwordService := crypto.PasswordService{}
	password, err := passwordService.Generate()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(password) != 16 {
		t.Errorf("expected password length 16, got %d", len(password))
	}
}

func TestHash(t *testing.T) {
	os.Setenv("ENCRYPT_SALT", "test_salt")
	defer os.Unsetenv("ENCRYPT_SALT")

	passwordService := crypto.PasswordService{}
	hash, err := passwordService.Hash("password123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Errorf("expected non-empty hash, got empty string")
	}
}

func TestHash_NoSalt(t *testing.T) {
	os.Unsetenv("ENCRYPT_SALT")

	passwordService := crypto.PasswordService{}
	hash, err := passwordService.Hash("password123")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "env variable ENCRYPT_SALT is not set" {
		t.Errorf("expected error 'env variable ENCRYPT_SALT is not set', got %v", err)
	}
	if hash != "" {
		t.Errorf("expected empty hash, got %v", hash)
	}
}

func TestVerifyPassword(t *testing.T) {
	os.Setenv("ENCRYPT_SALT", "test_salt")
	defer os.Unsetenv("ENCRYPT_SALT")

	passwordService := crypto.PasswordService{}
	hash, err := passwordService.Hash("password123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	valid, err := passwordService.Verify("password123", hash)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !valid {
		t.Errorf("expected valid password, got invalid")
	}
}

func TestVerify_InvalidHash(t *testing.T) {
	os.Setenv("ENCRYPT_SALT", "test_salt")
	defer os.Unsetenv("ENCRYPT_SALT")

	passwordService := crypto.PasswordService{}
	invalidHash := "invalidhash"

	valid, err := passwordService.Verify("password123", invalidHash)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if valid {
		t.Errorf("expected invalid password, got valid")
	}
}

func TestVerify_NoSalt(t *testing.T) {
	os.Unsetenv("ENCRYPT_SALT")

	passwordService := crypto.PasswordService{}
	hash, err := passwordService.Hash("password123")
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	valid, err := passwordService.Verify("password123", hash)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "env variable ENCRYPT_SALT is not set" {
		t.Errorf("expected error 'env variable ENCRYPT_SALT is not set', got %v", err)
	}
	if valid {
		t.Errorf("expected invalid password, got valid")
	}
}

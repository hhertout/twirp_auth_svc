package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"math/big"
	"os"

	"golang.org/x/crypto/argon2"
)

// PasswordServiceInterface defines the methods required for managing passwords.
// It includes methods for generating, hashing, and verifying passwords.
type PasswordServiceInterface interface {
	// Generate generates a random password of 16 characters using a predefined character set.
	// Returns the generated password and an error if any occurs.
	Generate() (string, error)

	// Hash generates a secure hash for a given password using the Argon2 function.
	// Uses an environment variable "ENCRYPT_SALT" as the salt.
	// Returns the base64 encoded hash and an error if any occurs.
	Hash(password string) (string, error)

	// Verify checks if a given password matches a hash using the Argon2 function.
	// Uses an environment variable "ENCRYPT_SALT" as the salt.
	// Returns a boolean indicating if the password is valid and an error if any occurs.
	Verify(password string, hash string) (bool, error)
}

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// Generate generates a random password of 16 characters using a predefined character set.
// Returns the generated password and an error if any occurs.
func (p *PasswordService) Generate() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*=+?/"

	passwordLength := 16
	charsetLen := big.NewInt(int64(len(charset)))

	password := make([]byte, passwordLength)
	for i := 0; i < passwordLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

func (p *PasswordService) Hash(password string) (string, error) {
	salt := os.Getenv("ENCRYPT_SALT")
	if salt == "" {
		return "", errors.New("env variable ENCRYPT_SALT is not set")
	}
	key := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	return base64.RawStdEncoding.EncodeToString(key), nil
}

func (p *PasswordService) Verify(password string, hash string) (bool, error) {
	salt := os.Getenv("ENCRYPT_SALT")
	if salt == "" {
		return false, errors.New("env variable ENCRYPT_SALT is not set")
	}

	decodeHash, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	hashToCompare := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	res := subtle.ConstantTimeCompare(hashToCompare, decodeHash)

	return res == 1, nil
}

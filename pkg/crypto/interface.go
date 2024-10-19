package crypto

import "github.com/golang-jwt/jwt/v5"

// JWTServiceInterface defines the methods required for managing JWT tokens.
// It includes methods for generating and verifying JWT tokens.
type JWTServiceInterface interface {
	// Generate creates a JWT token for a given user.
	// Uses an environment variable "JWT_SECRET" as the secret key.
	// The token expires in 20 days from the time of generation.
	//
	// Parameters:
	// - user: the identifier for the user for whom the token is being generated.
	//
	// Returns:
	// - The signed JWT token.
	// - An error if any occurs during the token generation.
	Generate(user string) (string, error)

	// Verify checks if a given JWT token is valid.
	// Uses an environment variable "JWT_SECRET" as the secret key.
	//
	// Parameters:
	// - tokenString: the JWT token to be verified.
	//
	// Returns:
	// - A boolean indicating if the token is valid.
	// - The token claims if the token is valid.
	// - An error if any occurs during the token verification.
	Verify(tokenString string) (bool, jwt.RegisteredClaims, error)
}

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

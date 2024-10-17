package server

import (
	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/pkg/auth"
	"github.com/hhertout/twirp_auth/pkg/crypto"
	"go.uber.org/zap"
)

// Server implements the different servers
type AuthenticationServer struct {
	Logger          *zap.Logger
	UserRepository  *repository.UserRepository
	JwtService      crypto.JWTServiceInterface
	PasswordService crypto.PasswordServiceInterface
}

// UserServer implements the different servers
type UserServer struct {
	Logger          *zap.Logger
	UserRepository  *repository.UserRepository
	JwtService      crypto.JWTServiceInterface
	PasswordService crypto.PasswordServiceInterface
	AuthManager     auth.AuthManagerInterface
}

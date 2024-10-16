package server

import (
	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/pkg/auth"
	"github.com/hhertout/twirp_auth/pkg/tools"
	"go.uber.org/zap"
)

// Server implements the different servers
type AuthenticationServer struct {
	UserRepository *repository.UserRepository
	Logger         *zap.Logger
}

// UserServer implements the different servers
type UserServer struct {
	UserRepository  *repository.UserRepository
	PasswordService tools.PasswordServiceInterface
	AuthManager     auth.AuthManagerInterface
	Logger          *zap.Logger
}

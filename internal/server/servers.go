package server

import "github.com/hhertout/twirp_auth/internal/repository"

// Server implements the different servers
type AuthenticationServer struct {
	UserRepository *repository.UserRepository
}

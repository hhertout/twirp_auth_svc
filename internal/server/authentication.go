package server

import (
	"context"

	"github.com/hhertout/twirp_auth/protobuf"
)

// AuthenticationServer implements the different servers
//
// @route /api/auth.AuthenticationService/Login
func (s *AuthenticationServer) Login(ctx context.Context, creds *protobuf.LoginRequest) (*protobuf.LoginResponse, error) {
	return &protobuf.LoginResponse{Token: "token"}, nil
}

// AuthenticationServer implements the different servers
//
// @route /api/auth.AuthenticationService/Register
func (s *AuthenticationServer) Register(ctx context.Context, creds *protobuf.RegisterRequest) (*protobuf.RegisterResponse, error) {
	return &protobuf.RegisterResponse{Token: "token"}, nil
}

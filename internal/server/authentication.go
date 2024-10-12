package server

import (
	"context"

	"github.com/hhertout/twirp_auth/internal/services"
	"github.com/hhertout/twirp_auth/pkg/tools"
	"github.com/hhertout/twirp_auth/protobuf"
	"github.com/twitchtv/twirp"
)

// AuthenticationServer implements the different servers
//
// @route /api/auth.AuthenticationService/Login
func (s *AuthenticationServer) Login(ctx context.Context, creds *protobuf.LoginRequest) (*protobuf.LoginResponse, error) {
	if err := services.CheckCredentials(creds.Username, creds.Password); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	user, err := s.UserRepository.FindOneByEmail(creds.Username)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	if user.Email != creds.Username {
		return nil, twirp.NotFound.Error("User not found")
	}

	match, err := tools.NewPasswordService().Verify(creds.Password, user.Password)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	if !match {
		return nil, twirp.Unauthenticated.Error("Invalid credentials")
	}

	token, err := tools.NewJWTService().Generate(creds.Username)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &protobuf.LoginResponse{Token: token}, nil
}

// AuthenticationServer implements the different servers
//
// @route /api/auth.AuthenticationService/Register
func (s *AuthenticationServer) Register(ctx context.Context, creds *protobuf.RegisterRequest) (*protobuf.RegisterResponse, error) {
	if err := services.CheckCredentials(creds.Username, creds.Password); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	user, err := s.UserRepository.FindOneByEmail(creds.Username)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	if user.Email == creds.Username {
		return nil, twirp.AlreadyExists.Error("User already exists")
	}

	hash, err := tools.NewPasswordService().Hash(creds.Password)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	_, err = s.UserRepository.Create(creds.Username, hash)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	token, err := tools.NewJWTService().Generate(creds.Username)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &protobuf.RegisterResponse{Token: token, Username: creds.Username}, nil
}

func (s *AuthenticationServer) CheckToken(ctx context.Context, req *protobuf.CheckTokenRequest) (*protobuf.CheckTokenResponse, error) {
	token := req.GetToken()
	if token == "" {
		return nil, twirp.InvalidArgument.Error("Token is empty")
	}

	valid, claims, err := tools.NewJWTService().Verify(token)
	if err != nil {
		return nil, twirp.Unauthenticated.Error(err.Error())
	}

	if !valid {
		return nil, twirp.Unauthenticated.Error("Invalid token")
	}

	return &protobuf.CheckTokenResponse{Username: claims.Issuer}, nil
}

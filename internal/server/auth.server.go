package server

import (
	"context"

	"github.com/hhertout/twirp_auth/internal/services"
	"github.com/hhertout/twirp_auth/pkg/tools"
	"github.com/hhertout/twirp_auth/protobuf/proto_auth"
	"github.com/twitchtv/twirp"
)

// AuthenticationServer implements the different servers
//
// @route /api/auth.AuthenticationService/Login
func (s *AuthenticationServer) Login(ctx context.Context, creds *proto_auth.LoginRequest) (*proto_auth.LoginResponse, error) {
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

	return &proto_auth.LoginResponse{Token: token}, nil
}

func (s *AuthenticationServer) CheckToken(ctx context.Context, req *proto_auth.CheckTokenRequest) (*proto_auth.CheckTokenResponse, error) {
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

	return &proto_auth.CheckTokenResponse{Username: claims.Issuer}, nil
}

package server

import (
	"context"
	"strings"

	"github.com/hhertout/twirp_auth/internal/services"
	"github.com/hhertout/twirp_auth/pkg/auth/role"
	"github.com/hhertout/twirp_auth/protobuf/proto_user"
	"github.com/twitchtv/twirp"
)

// UserServer implements the different servers
//
// @route /api/user.UserService/Register
func (u *UserServer) Register(ctx context.Context, req *proto_user.RegisterRequest) (*proto_user.RegisterResponse, error) {
	// Check user data
	// this is to move in a dedicated service
	if req.Username == "" || !strings.Contains(req.Username, "@") {
		u.Logger.Sugar().Error("Username is not a valid email")
		return nil, twirp.InvalidArgumentError(req.Username, "Username is not a valid email")
	}
	if req.Password == "" {
		u.Logger.Sugar().Error("Password is empty")
		return nil, twirp.InvalidArgumentError(req.Password, "Password are required")
	}

	if err := services.CheckCredentials(req.Username, req.Password); err != nil {
		u.Logger.Sugar().Error("Error during the check of the credentials", err)
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	user, err := u.UserRepository.FindOneByEmailInAll(req.Username)
	if err != nil {
		u.Logger.Sugar().Error("Error during the search of the user", err)
		return nil, twirp.InternalErrorWith(err)
	}

	if user.DeletedAt != "" {
		u.Logger.Sugar().Error("User is banned", user.Email)
		return nil, twirp.PermissionDenied.Error("User is banned")
	}

	if user.Email == req.Username {
		u.Logger.Sugar().Error("User already exists", user.Email)
		return nil, twirp.AlreadyExists.Error("User already exists")
	}

	hash, err := u.PasswordService.Hash(req.Password)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	_, err = u.UserRepository.Create(req.Username, hash, []string{"USER"})
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	token, err := u.JwtService.Generate(req.Username)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &proto_user.RegisterResponse{Token: token, Username: req.Username}, nil
}

// UserServer implements the different servers
//
// @route /api/user.UserService/Login
func (u *UserServer) Ban(ctx context.Context, req *proto_user.BanRequest) (*proto_user.BanResponse, error) {
	user, err := u.AuthManager.AllowAccessWithRole(ctx, []role.ROLE{role.ROLE_ADMIN})
	if err != nil {
		u.Logger.Sugar().Error("Error during the check of the credentials", err)
		return nil, twirp.PermissionDenied.Error(err.Error())
	}

	if req.Username == "" {
		u.Logger.Sugar().Error("Username is empty")
		return nil, twirp.InvalidArgument.Error("Username is empty")
	}

	if user.Email != req.Username {
		u.Logger.Sugar().Error("User not found")
		return nil, twirp.NotFound.Error("User not found")
	}

	_, err = u.UserRepository.SoftDelete(user.Email)
	if err != nil {
		u.Logger.Sugar().Error("Error during the soft delete of the user", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &proto_user.BanResponse{Success: true}, nil
}

// UserServer implements the different servers
//
// @route /api/user.UserService/Unban
func (u *UserServer) Unban(ctx context.Context, req *proto_user.UnbanRequest) (*proto_user.UnbanResponse, error) {
	user, err := u.AuthManager.AllowAccessWithRole(ctx, []role.ROLE{role.ROLE_ADMIN})
	if err != nil {
		u.Logger.Sugar().Error("Error during the check of the credentials", err)
		return nil, twirp.PermissionDenied.Error(err.Error())
	}

	if req.Username == "" {
		u.Logger.Sugar().Error("Username is empty")
		return nil, twirp.InvalidArgument.Error("Username is empty")
	}

	if user.Email != req.Username {
		u.Logger.Sugar().Error("User not found")
		return nil, twirp.NotFound.Error("User not found")
	}

	_, err = u.UserRepository.RemoveSoftDelete(user.Email)
	if err != nil {
		u.Logger.Sugar().Error("Error during the soft delete of the user", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &proto_user.UnbanResponse{Success: true}, nil
}

// UserServer implements the different servers
//
// @route /api/user.UserService/Update
func (u *UserServer) Delete(ctx context.Context, req *proto_user.DeleteRequest) (*proto_user.DeleteResponse, error) {
	user, err := u.AuthManager.AllowAccessWithRole(ctx, []role.ROLE{role.ROLE_ADMIN})
	if err != nil {
		u.Logger.Sugar().Error("Error during the check of the credentials", err)
		return nil, twirp.PermissionDenied.Error(err.Error())
	}

	if req.Username == "" {
		u.Logger.Sugar().Error("Username is empty")
		return nil, twirp.InvalidArgument.Error("Username is empty")
	}

	if user.Email != req.Username {
		u.Logger.Sugar().Error("User not found")
		return nil, twirp.NotFound.Error("User not found")
	}

	_, err = u.UserRepository.HardDelete(user.Id)
	if err != nil {
		u.Logger.Sugar().Error("Error during the hard delete of the user", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &proto_user.DeleteResponse{Success: true}, nil
}

// UserServer implements the different servers
//
// @route /api/user.UserService/Update
func (u *UserServer) UpdatePassword(ctx context.Context, req *proto_user.UpdatePasswordRequest) (*proto_user.UpdatePasswordResponse, error) {
	user, err := u.AuthManager.AllowAccessWithRole(ctx, []role.ROLE{})
	if err != nil {
		u.Logger.Sugar().Error("Error during the check of the credentials", err)
		return nil, twirp.PermissionDenied.Error(err.Error())
	}

	if req.Username == "" {
		u.Logger.Sugar().Error("Username is empty")
		return nil, twirp.InvalidArgument.Error("Username is empty")
	}

	if req.NewPassword == "" {
		u.Logger.Sugar().Error("Password is empty")
		return nil, twirp.InvalidArgument.Error("Password is empty")
	}

	if user.Email != req.Username {
		u.Logger.Sugar().Error("User not found")
		return nil, twirp.NotFound.Error("User not found")
	}

	valid, err := u.PasswordService.Verify(req.OldPassword, user.Password)
	if err != nil {
		u.Logger.Sugar().Error("Error during the verification of the password", err)
		return nil, twirp.InternalErrorWith(err)
	}

	if !valid {
		return nil, twirp.Unauthenticated.Error("Invalid password")
	}

	hash, err := u.PasswordService.Hash(req.NewPassword)
	if err != nil {
		u.Logger.Sugar().Error("Error during the hashing of the password", err)
		return nil, twirp.InternalErrorWith(err)
	}

	_, err = u.UserRepository.UpdatePassword(user.Id, hash)
	if err != nil {
		u.Logger.Sugar().Error("Error during the update of the password", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &proto_user.UpdatePasswordResponse{Success: true}, nil
}

// UserServer implements the different servers
//
// @route /api/user.UserService/Update
func (u *UserServer) UpdateEmail(ctx context.Context, req *proto_user.UpdateEmailRequest) (*proto_user.UpdateEmailResponse, error) {
	user, err := u.AuthManager.AllowAccessWithRole(ctx, []role.ROLE{})
	if err != nil {
		u.Logger.Sugar().Error("Error during the check of the credentials", err)
		return nil, twirp.PermissionDenied.Error(err.Error())
	}

	if req.NewEmail == "" {
		u.Logger.Sugar().Error("New email is empty")
		return nil, twirp.InvalidArgument.Error("Username is empty")
	}

	if req.OldEmail == "" {
		u.Logger.Sugar().Error("Old email is empty")
		return nil, twirp.InvalidArgument.Error("New email is empty")
	}

	if user.Email == req.NewEmail {
		u.Logger.Sugar().Error("User already exists", user.Email)
		return nil, twirp.NotFound.Error("User already exists")
	}

	user, err = u.UserRepository.FindOneByEmail(req.OldEmail)
	if err != nil {
		u.Logger.Sugar().Error("Error during the search of the user", err)
		return nil, twirp.InternalErrorWith(err)
	}

	if user.Email != req.OldEmail {
		u.Logger.Sugar().Error("User not found")
		return nil, twirp.NotFound.Error("User not found")
	}

	_, err = u.UserRepository.UpdateEmail(req.OldEmail, req.NewEmail)
	if err != nil {
		u.Logger.Sugar().Error("Error during the update of the email", err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &proto_user.UpdateEmailResponse{Success: true}, nil
}

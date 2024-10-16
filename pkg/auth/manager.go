package auth

import (
	"context"
	"errors"

	"github.com/hhertout/twirp_auth/internal/hooks"
	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/pkg/auth/role"
	"github.com/hhertout/twirp_auth/pkg/tools"
)

// AuthManagerInterface defines the methods required for managing authentication and authorization.
// It includes methods for restricting and allowing access based on user roles.
type AuthManagerInterface interface {
	// RestrictAccessWithRole restricts access to a user based on their role.
	// It verifies the JWT token from the context and checks if the user has one of the required roles.
	// If the token is missing, invalid, or the user does not have the required role, it returns an error.
	// It return an error if the user have the role specified in the slice or the user otherwise.
	//
	// Parameters:
	// - ctx: the context containing the JWT token.
	// - roles: a slice of roles that are required for access.
	//
	// Returns:
	// - The user if they have the required role.
	// - An error if the token is missing, invalid, or the user does not have the required role.
	RestrictAccessWithRole(ctx context.Context, roles []role.ROLE) (repository.User, error)

	// AllowAccessWithRole allows access to a user based on their role.
	// It verifies the JWT token from the context and checks if the user has one of the required roles.
	// If the token is missing, invalid, or the user does not have the required role, it returns an error.
	// It return an error if the user does not have the role specified in the slice, or the user otherwise.
	//
	// Parameters:
	// - ctx: the context containing the JWT token.
	// - roles: a slice of roles that are required for access.
	//
	// Returns:
	// - The user if they have the required role.
	// - An error if the token is missing, invalid, or the user does not have the required role.
	AllowAccessWithRole(ctx context.Context, roles []role.ROLE) (repository.User, error)
}

type AuthManager struct {
	Repository *repository.UserRepository
	JWTManager tools.JWTServiceInterface
}

func NewAuthManager(r *repository.UserRepository) *AuthManager {
	return &AuthManager{
		Repository: r,
		JWTManager: tools.NewJWTService(),
	}
}

func (am *AuthManager) RestrictAccessWithRole(ctx context.Context, roles []role.ROLE) (repository.User, error) {
	token, _ := ctx.Value(hooks.ServerContextKey("Authorization")).(string)
	if token == "" {
		return repository.User{}, errors.New("token is missing")
	}

	isValid, claims, err := am.JWTManager.Verify(token)
	if err != nil {
		return repository.User{}, err
	}
	if !isValid {
		return repository.User{}, errors.New("invalid token")
	}

	user, err := am.Repository.FindOneByEmail(claims.Issuer)
	if err != nil {
		return repository.User{}, err
	}

	if len(roles) > 0 {
		for _, r := range user.Role {
			if role.Contains(roles, r) {
				return repository.User{}, errors.New("user does not have the required role")
			}
		}

		return user, nil
	} else {
		return user, nil
	}
}

func (am *AuthManager) AllowAccessWithRole(ctx context.Context, roles []role.ROLE) (repository.User, error) {
	token, _ := ctx.Value(hooks.ServerContextKey("Authorization")).(string)
	if token == "" {
		return repository.User{}, errors.New("token is missing")
	}

	isValid, claims, err := am.JWTManager.Verify(token)

	if err != nil {
		return repository.User{}, err
	}
	if !isValid {
		return repository.User{}, errors.New("invalid token")
	}

	user, err := am.Repository.FindOneByEmail(claims.Issuer)
	if err != nil {
		return repository.User{}, err
	}

	if len(roles) > 0 {
		for _, r := range user.Role {
			if role.Contains(roles, r) {
				return user, nil
			}
		}
	}

	return repository.User{}, errors.New("user does not have the required role")
}

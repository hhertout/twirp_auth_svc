package auth

import (
	"context"

	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/pkg/auth/role"
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

type AuthDataLayerInterface interface {
	FindOneByEmail(email string) (repository.User, error)
}

package auth

import (
	"context"
	"errors"

	"github.com/hhertout/twirp_auth/internal/hooks"
	"github.com/hhertout/twirp_auth/pkg/auth/role"
	"github.com/hhertout/twirp_auth/pkg/crypto"
	"github.com/hhertout/twirp_auth/pkg/dto"
)

type AuthManager struct {
	Dal        AuthDataLayerInterface
	JWTManager crypto.JWTServiceInterface
}

func NewAuthManager(r AuthDataLayerInterface) *AuthManager {
	return &AuthManager{
		Dal:        r,
		JWTManager: crypto.NewJWTService(),
	}
}

func (am *AuthManager) RestrictAccessWithRole(ctx context.Context, roles []role.ROLE) (dto.User, error) {
	token, _ := ctx.Value(hooks.ServerContextKey("Authorization")).(string)
	if token == "" {
		return dto.User{}, errors.New("token is missing")
	}

	isValid, claims, err := am.JWTManager.Verify(token)
	if err != nil {
		return dto.User{}, err
	}
	if !isValid {
		return dto.User{}, errors.New("invalid token")
	}

	user, err := am.Dal.FindOneByEmail(claims.Issuer)
	if err != nil {
		return dto.User{}, err
	}

	if len(roles) > 0 {
		for _, r := range user.Role {
			if role.Contains(roles, r) {
				return dto.User{}, errors.New("user does not have the required role")
			}
		}

		return user, nil
	} else {
		return user, nil
	}
}

func (am *AuthManager) AllowAccessWithRole(ctx context.Context, roles []role.ROLE) (dto.User, error) {
	token, _ := ctx.Value(hooks.ServerContextKey("Authorization")).(string)
	if token == "" {
		return dto.User{}, errors.New("token is missing")
	}

	isValid, claims, err := am.JWTManager.Verify(token)

	if err != nil {
		return dto.User{}, err
	}
	if !isValid {
		return dto.User{}, errors.New("invalid token")
	}

	user, err := am.Dal.FindOneByEmail(claims.Issuer)
	if err != nil {
		return dto.User{}, err
	}

	if len(roles) > 0 {
		for _, r := range user.Role {
			if role.Contains(roles, r) {
				return user, nil
			}
		}
	}

	return dto.User{}, errors.New("user does not have the required role")
}

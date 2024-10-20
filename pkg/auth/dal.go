package auth

import (
	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/pkg/dto"
)

type AuthDataLayer struct {
	repository *repository.UserRepository
}

func NewAuthDataLayer(r *repository.UserRepository) *AuthDataLayer {
	return &AuthDataLayer{r}
}

func (dl *AuthDataLayer) FindOneByEmail(email string) (dto.User, error) {
	return dl.repository.FindOneByEmail(email)
}

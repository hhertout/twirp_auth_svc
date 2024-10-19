package auth

import "github.com/hhertout/twirp_auth/internal/repository"

type AuthDataLayerInterface interface {
	FindOneByEmail(email string) (repository.User, error)
}

type AuthDataLayer struct {
	repository *repository.UserRepository
}

func NewAuthDataLayer(r *repository.UserRepository) *AuthDataLayer {
	return &AuthDataLayer{r}
}

func (dl *AuthDataLayer) FindOneByEmail(email string) (repository.User, error) {
	return dl.repository.FindOneByEmail(email)
}

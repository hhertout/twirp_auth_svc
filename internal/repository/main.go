package repository

import (
	"database/sql"

	"github.com/hhertout/twirp_auth/pkg/database"
)

type UserRepository struct {
	dbPool *sql.DB
}

// NewRepository creates a new instance of Repository.
// If a custom database source is provided, it uses that source.
// Otherwise, it connects to the default database.
//
// Parameters:
// - customSource: a pointer to a sql.DB instance representing a custom database source.
//
// Returns:
// - A pointer to the newly created Repository.
// - An error if any occurs during the connection to the default database.
func NewUserRepository(customSource *sql.DB) (*UserRepository, error) {
	if customSource != nil {
		return &UserRepository{
			customSource,
		}, nil
	} else {
		dbService, err := database.Connect()
		if err != nil {
			return nil, err
		}

		return &UserRepository{
			dbService.DbPool,
		}, nil
	}
}

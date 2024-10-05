package services

import (
	"errors"
)

func CheckCredentials(username string, password string) error {
	if username == "" {
		return errors.New("username is required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	return nil
}

package repository

import (
	"database/sql"

	"github.com/hhertout/twirp_auth/pkg/database"
	"github.com/hhertout/twirp_auth/pkg/dto"
)

type UserRepository struct {
	dbPool *sql.DB
}

// NewRepository creates a new instance of Repository.
// If a custom database source is provided, it uses that source.
// Otherwise, it connects to the default database.
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

func (r UserRepository) Create(email string, password string, role []string) (int, error) {
	res, err := r.dbPool.Exec(`
		INSERT INTO "user" (email, password, role) 
		VALUES ($1, $2, $3)
	`, email, password, role)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (r UserRepository) FindOneByEmail(email string) (dto.User, error) {
	var user dto.User
	rows, err := r.dbPool.Query(`
		SELECT id, uuid, email, password 
		FROM "user" 
		WHERE email=$1 AND deleted_at is null 
		LIMIT 1
	`, email)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Uuuid, &user.Email, &user.Password)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (r UserRepository) FindCompleteOneByEmail(email string) (dto.CompleteUser, error) {
	var user dto.CompleteUser
	rows, err := r.dbPool.Query(`
		SELECT id, uuid, email, password, deleted_at, created_at, updated_at 
		FROM "user" 
		WHERE email=$1
		LIMIT 1
	`, email)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.Id,
			&user.Uuid,
			&user.Email,
			&user.Password,
			&user.DeletedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (r UserRepository) UpdatePassword(id string, password string) (int, error) {
	res, err := r.dbPool.Exec(`
		UPDATE "user" 
		SET password=$1 
		where id=$2
	`, password, id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (r UserRepository) UpdateEmail(oldEmail string, newEmail string) (int, error) {
	res, err := r.dbPool.Exec(`
		UPDATE "user" 
		SET email=$1 
		where email=$2
	`, newEmail, oldEmail)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (r UserRepository) SoftDelete(email string) (int, error) {
	res, err := r.dbPool.Exec(`
		UPDATE "user"
		SET deleted_at=NOW()
		WHERE email=$1;
	`, email)

	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (r UserRepository) RemoveSoftDelete(email string) (int, error) {
	res, err := r.dbPool.Exec(`
		UPDATE "user"
		SET deleted_at=NULL
		WHERE email=$1;
	`, email)

	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (r UserRepository) HardDelete(email string) (int, error) {
	res, err := r.dbPool.Exec(`
		DELETE FROM "user"
		WHERE email=$1;
	`, email)

	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

package repository

type User struct {
	Id       string   `db:"id"`
	Uuuid    string   `db:"uuid"`
	Email    string   `db:"email"`
	Password string   `db:"password"`
	Role     []string `db:"role"`
}

type CompleteUser struct {
	Id        string   `db:"id"`
	Uuid      string   `db:"uuid"`
	Email     string   `db:"email"`
	Password  string   `db:"password"`
	Role      []string `db:"role"`
	CreatedAt string   `db:"created_at"`
	UpdatedAt string   `db:"updated_at"`
	DeletedAt string   `db:"deleted_at"`
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

func (r UserRepository) FindOneByEmail(email string) (User, error) {
	var user User
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

func (r UserRepository) FindOneByEmailInAll(email string) (CompleteUser, error) {
	var user CompleteUser
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

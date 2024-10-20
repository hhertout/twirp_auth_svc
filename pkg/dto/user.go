package dto

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

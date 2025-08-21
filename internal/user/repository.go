package user

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(user *User) error {
	_, err := repo.db.Exec(`INSERT INTO users (id, email, password_hash, created_at, updated_at)
         VALUES ($1, $2, $3, NOW(), NOW())`,
		 user.ID, user.Email, user.PasswordHash,
	)

	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return ErrEmailTaken
			}
		}
		return err
	}
	
	return nil
}

func (repo *Repository) GetByEmail(email string) (*User, error) {
	var user User
	err := repo.db.QueryRow(
		`SELECT id,email,password_hash,created_at,updated_at,deleted_at
       FROM users
      WHERE email=$1 AND deleted_at IS NULL`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
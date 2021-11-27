package mysql

import (
	"database/sql"
	"fmt"

	"github.com/bazeeko/investor-social-network/domain"
)

type mysqlUserRepository struct {
	*sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.MysqlUserRepository {
	return &mysqlUserRepository{db}
}

func (r *mysqlUserRepository) Add(u domain.User) error {
	_, err := r.Exec(`INSERT users (username, password, rating, created_at) VALUES (?, ?, ?, ?)`,
		u.Username, u.Password, u.Rating, u.CreatedAt)

	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}

	return nil
}

func (r *mysqlUserRepository) GetById(id int) (domain.User, error) {
	u := domain.User{}

	err := r.QueryRow(`SELECT id, username, rating, created_at FROM users WHERE id=?`, id).
		Scan(&u.ID, &u.Username, &u.Rating, &u.CreatedAt)

	if err != nil {
		return domain.User{}, fmt.Errorf("GetById: %w", err)
	}

	return u, nil
}

func (r *mysqlUserRepository) GetByUsername(username string) (domain.User, error) {
	u := domain.User{}

	err := r.QueryRow(`SELECT id, username, rating, created_at FROM users WHERE username=?`, username).
		Scan(&u.ID, &u.Username, &u.Rating, &u.CreatedAt)

	if err != nil {
		return domain.User{}, fmt.Errorf("GetById: %w", err)
	}

	return u, nil
}

func (r *mysqlUserRepository) GetUserCredentials(username string) (string, string, error) {
	var password string

	err := r.QueryRow(`SELECT password FROM users WHERE username=?`, username).
		Scan(&password)

	if err != nil {
		return "", "", fmt.Errorf("GetByUsername: %w", err)
	}

	return username, password, nil
}

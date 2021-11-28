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

func (r *mysqlUserRepository) AddUserToFavourites(userID int, favUserID int) error {
	_, err := r.Exec(`INSERT favourite_people (
		user_id, favourite_user_id)
		VALUES (?, ?)`,
		userID, favUserID)

	if err != nil {
		return fmt.Errorf("AddUserToFavourites: %w", err)
	}

	return nil
}

func (r *mysqlUserRepository) DeleteUserFromFavourites(userID int, favUserID int) error {
	_, err := r.Exec(`DELETE FROM favourite_people
		WHERE user_id=? AND favourite_user_id=?`,
		userID, favUserID)

	if err != nil {
		return fmt.Errorf("DeleteUserFromFavourites: %w", err)
	}

	return nil
}

func (r *mysqlUserRepository) GetFavouriteUsers(userID int) ([]int, error) {
	users := make([]int, 0)

	rows, err := r.Query(`SELECT favourite_user_id FROM favourite_people
	WHERE user_id=?`, userID)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return []int{}, err
		}

		users = append(users, id)
	}

	if err != nil {
		return []int{}, fmt.Errorf("GetFavouriteUsers: %w", err)
	}

	return users, nil
}

func (r *mysqlUserRepository) AddLikeToUser(userID int, likedUserID int) error {
	_, err := r.Exec(`INSERT user_likes (
		user_id, liked_user_id)
		VALUES (?, ?)`, userID, likedUserID)

	if err != nil {
		return fmt.Errorf("AddLikeToUser: %w", err)
	}

	r.Exec(`UPDATE users SET rating=rating+1 WHERE id=?`, likedUserID)

	return nil
}

func (r *mysqlUserRepository) DeleteLikeFromUser(userID int, likedUserID int) error {
	_, err := r.Exec(`DELETE FROM user_likes
	WHERE user_id=? AND liked_user_id=?`, userID, likedUserID)

	if err != nil {
		return fmt.Errorf("DeleteLikeFromUser: %w", err)
	}

	r.Exec(`UPDATE users SET rating=rating-1 WHERE id=?`, likedUserID)

	return nil
}

// `DELETE FROM favourite_stocks
// WHERE user_id=? AND stock_symbol=?`,

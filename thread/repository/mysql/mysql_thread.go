package mysql

import (
	"database/sql"
	"fmt"

	"github.com/bazeeko/investor-social-network/domain"
)

type mysqlThreadRepository struct {
	*sql.DB
}

func NewMysqlThreadRepository(db *sql.DB) domain.MysqlThreadRepository {
	return &mysqlThreadRepository{db}
}

func (r *mysqlThreadRepository) CreateThread(userID int, t domain.Thread) error {
	_, err := r.Exec(`INSERT threads (
		user_id, topic, body, image_url, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		userID, t.Topic, t.Body, t.ImageURL, t.CreatedAt)

	if err != nil {
		return fmt.Errorf("CreateThread: %w", err)
	}

	return nil
}

func (r *mysqlThreadRepository) GetThreadByID(userID int, threadID int) (domain.Thread, error) {
	t := domain.Thread{}

	err := r.QueryRow(`SELECT
		id, topic, body, image_url, created_at
		FROM threads
		WHERE user_id=? AND thread_id=?`, userID, threadID).
		Scan(&t.ID, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)

	if err != nil {
		return domain.Thread{}, fmt.Errorf("GetThreadByID: %w", err)
	}

	return t, nil
}

// func (r *mysqlThreadRepository) GetThreads(userID int, threadID int) (domain.Thread, error) {
// 	t := domain.Thread{}

// 	err := r.QueryRow(`SELECT
// 		id, topic, body, image_url, created_at
// 		FROM threads
// 		WHERE user_id=? AND thread_id=?`, userID, threadID).
// 		Scan(t.ID, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)

// 	if err != nil {
// 		return domain.Thread{}, fmt.Errorf("GetThreadByID: %w", err)
// 	}

// 	return t, nil
// }

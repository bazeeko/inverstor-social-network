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

func (r *mysqlThreadRepository) GetThreadByID(threadID int) (domain.Thread, error) {
	t := domain.Thread{}

	err := r.QueryRow(`SELECT
		id, hashtag, topic, body, image_url, created_at
		FROM threads
		WHERE thread_id=?`, threadID).
		Scan(&t.ID, &t.Hashtag, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)

	if err != nil {
		return domain.Thread{}, fmt.Errorf("GetThreadByID: %w", err)
	}

	return t, nil
}

func (r *mysqlThreadRepository) DeleteThreadByID(threadID int) error {
	_, err := r.Exec(`DELETE FROM threads
	WHERE id=?`, threadID)
	if err != nil {
		return fmt.Errorf("DeleteThreadByID: %w", err)
	}

	return nil
}

func (r *mysqlThreadRepository) GetUserThreads(userID int) ([]domain.Thread, error) {
	threads := []domain.Thread{}

	rows, err := r.Query(`SELECT
	id, hashtag, topic, body, image_url, created_at
	FROM threads
	WHERE user_id=?`, userID)

	if err != nil {
		return []domain.Thread{}, fmt.Errorf("GetUserThreads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := domain.Thread{}

		err := rows.Scan(&t.ID, &t.Hashtag, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)
		if err != nil {
			return []domain.Thread{}, fmt.Errorf("GetUserThreads: %w", err)
		}

		threads = append(threads, t)
	}

	return threads, nil
}

func (r *mysqlThreadRepository) GetThreadsByHashtag(hashtag string) ([]domain.Thread, error) {
	threads := []domain.Thread{}

	rows, err := r.Query(`SELECT
	id, hashtag, topic, body, image_url, created_at
	FROM threads
	WHERE hashtag=?`, hashtag)

	if err != nil {
		return []domain.Thread{}, fmt.Errorf("GetThreadsByHashtag: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := domain.Thread{}

		err := rows.Scan(&t.ID, &t.Hashtag, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)
		if err != nil {
			return []domain.Thread{}, fmt.Errorf("GetThreadsByHashtag: %w", err)
		}

		threads = append(threads, t)
	}

	return threads, nil
}

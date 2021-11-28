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
		user_id, hashtag, topic, body, image_url, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		userID, t.Hashtag, t.Topic, t.Body, t.ImageURL, t.CreatedAt)

	if err != nil {
		return fmt.Errorf("CreateThread: %w", err)
	}

	return nil
}

func (r *mysqlThreadRepository) GetThreadByID(threadID int) (domain.Thread, error) {
	t := domain.Thread{}

	err := r.QueryRow(`SELECT
		id, user_id, hashtag, topic, body, image_url, created_at
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
	id, user_id, hashtag, topic, body, image_url, created_at
	FROM threads
	WHERE user_id=?`, userID)

	if err != nil {
		return []domain.Thread{}, fmt.Errorf("GetUserThreads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := domain.Thread{}

		err := rows.Scan(&t.ID, &t.UserID, &t.Hashtag, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)
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
	id, user_id, hashtag, topic, body, image_url, created_at
	FROM threads
	WHERE hashtag=?`, hashtag)

	if err != nil {
		return []domain.Thread{}, fmt.Errorf("GetThreadsByHashtag: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := domain.Thread{}

		err := rows.Scan(&t.ID, &t.UserID, &t.Hashtag, &t.Topic, &t.Body, &t.ImageURL, &t.CreatedAt)
		if err != nil {
			return []domain.Thread{}, fmt.Errorf("GetThreadsByHashtag: %w", err)
		}

		threads = append(threads, t)
	}

	return threads, nil
}

func (r *mysqlThreadRepository) CreateComment(c domain.Comment) error {
	_, err := r.Exec(`INSERT comments (
		user_id, thread_id, body, created_at) 
		VALUES (?, ?, ?, ?)`,
		c.UserID, c.ThreadID, c.Body, c.CreatedAt)

	if err != nil {
		return fmt.Errorf("CreateComment: %w", err)
	}

	return nil
}

func (r *mysqlThreadRepository) GetCommentsByThreadID(threadID int) ([]domain.Comment, error) {
	comments := []domain.Comment{}

	rows, err := r.Query(`SELECT
	id, user_id, thread_id, body, created_at
	FROM comments
	WHERE thread_id=?`, threadID)

	if err != nil {
		return []domain.Comment{}, fmt.Errorf("GetCommentsByThreadID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		c := domain.Comment{}

		err := rows.Scan(&c.ID, &c.UserID, &c.ThreadID, &c.Body, &c.CreatedAt)
		if err != nil {
			return []domain.Comment{}, fmt.Errorf("GetCommentsByThreadID: %w", err)
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func (r *mysqlThreadRepository) CreateSubComment(sc domain.SubComment) error {
	_, err := r.Exec(`INSERT sub_comments (
		user_id, thread_id, comment_id, body, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		sc.UserID, sc.ThreadID, sc.CommentID, sc.Body, sc.CreatedAt)

	if err != nil {
		return fmt.Errorf("CreateSubComment: %w", err)
	}

	return nil
}

func (r *mysqlThreadRepository) GetSubCommentsByCommentID(commentID int) ([]domain.SubComment, error) {
	subcomments := []domain.SubComment{}

	rows, err := r.Query(`SELECT
	id, user_id, thread_id, comment_id, body, created_at
	FROM sub_comments
	WHERE comment_id=?`, commentID)

	if err != nil {
		return []domain.SubComment{}, fmt.Errorf("GetSubCommentsByCommentID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		sc := domain.SubComment{}

		err := rows.Scan(&sc.ID, &sc.UserID, &sc.ThreadID, &sc.CommentID, &sc.Body, &sc.CreatedAt)
		if err != nil {
			return []domain.SubComment{}, fmt.Errorf("GetSubCommentsByCommentID: %w", err)
		}

		subcomments = append(subcomments, sc)
	}

	return subcomments, nil
}

func (r *mysqlThreadRepository) AddLikeToThread(userID int, threadID int) error {
	_, err := r.Exec(`INSERT thread_likes (
		user_id, thread_id)
		VALUES (?, ?)`, userID, threadID)

	if err != nil {
		return fmt.Errorf("AddLikeToThread: %w", err)
	}

	var uID int

	r.QueryRow(`SELECT user_id FROM threads WHERE id=?`, threadID).Scan(&uID)

	r.Exec(`UPDATE users SET rating=rating+0.1 WHERE id=?`, uID)

	return nil
}

func (r *mysqlThreadRepository) DeleteLikeFromThread(userID int, threadID int) error {
	_, err := r.Exec(`DELETE FROM thread_likes
	WHERE user_id=? AND thread_id=?`, userID, threadID)

	if err != nil {
		return fmt.Errorf("DeleteLikeFromThread: %w", err)
	}

	var uID int

	r.QueryRow(`SELECT user_id FROM threads WHERE id=?`, threadID).Scan(&uID)

	r.Exec(`UPDATE users SET rating=rating-0.1 WHERE id=?`, uID)

	return nil
}

// _, err = conn.Exec(`CREATE TABLE IF NOT EXISTS comments (
// 	id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
// 	user_id BIGINT NOT NULL,
// 	thread_id BIGINT NOT NULL,
// 	body TEXT,
// 	image_url TEXT,
// 	created_at VARCHAR(40) NOT NULL,
// 	PRIMARY KEY (id),
// 	FOREIGN KEY (user_id) REFERENCES users(id),
// 	FOREIGN KEY (thread_id) REFERENCES threads(id)
// );`)

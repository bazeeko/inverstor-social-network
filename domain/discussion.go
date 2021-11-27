package domain

type Thread struct {
	ID        int
	Hashtag   string
	Topic     string
	Body      string
	ImageURL  string
	CreatedAt string
}

type Comment struct {
	ID        int
	Body      string
	ImageURL  string
	CreatedAt string
}

type SubComment struct {
	ID        int
	CommentID int
	Body      string
	CreatedAt string
}

type MysqlThreadRepository interface {
	CreateThread(userID int, t Thread) error
	GetThreadByID(threadID int) (Thread, error)
	DeleteThreadByID(threadID int) error
}

type ThreadUsecase interface {
}

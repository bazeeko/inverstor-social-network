package domain

type Thread struct {
	ID        int
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
}

type ThreadUsecase interface {
}

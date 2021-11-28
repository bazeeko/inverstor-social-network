package domain

type Thread struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Hashtag   string    `json:"hashtag"`
	Topic     string    `json:"topic"`
	Body      string    `json:"body"`
	ImageURL  string    `json:"image_url"`
	Likes     int       `json:"likes"`
	CreatedAt string    `json:"created_at"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	ID          int          `json:"id"`
	UserID      int          `json:"user_id"`
	ThreadID    int          `json:"thread_id"`
	Body        string       `json:"body"`
	CreatedAt   string       `json:"created_at"`
	SubComments []SubComment `json:"sub_comments"`
}

type SubComment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ThreadID  int    `json:"thread_id"`
	CommentID int    `json:"comment_id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type MysqlThreadRepository interface {
	CreateThread(userID int, t Thread) error
	GetThreadByID(threadID int) (Thread, error)
	DeleteThreadByID(threadID int) error
	GetUserThreads(userID int) ([]Thread, error)
	GetThreadsByHashtag(hashtag string) ([]Thread, error)
	CreateComment(c Comment) error
	CreateSubComment(sc SubComment) error
	AddLikeToThread(userID int, threadID int) error
	DeleteLikeFromThread(userID int, threadID int) error
	GetCommentsByThreadID(threadID int) ([]Comment, error)
	GetSubCommentsByCommentID(commentID int) ([]SubComment, error)
	GetAmountOfLikes(threadID int) (int, error)
}

type ThreadUsecase interface {
	CreateThread(userID int, t Thread) error
	GetThreadByID(threadID int) (Thread, error)
	DeleteThreadByID(threadID int) error
	GetUserThreads(userID int) ([]Thread, error)
	GetThreadsByHashtag(hashtag string) ([]Thread, error)
	CreateComment(c Comment) error
	CreateSubComment(sc SubComment) error
	AddLikeToThread(userID int, threadID int) error
	DeleteLikeFromThread(userID int, threadID int) error
	GetCommentsByThreadID(threadID int) ([]Comment, error)
	GetSubCommentsByCommentID(commentID int) ([]SubComment, error)
}

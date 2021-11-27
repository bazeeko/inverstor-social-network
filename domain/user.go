package domain

type User struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Password  string  `json:"-"`
	Rating    float64 `json:"rating"`
	CreatedAt string  `json:"created_at"`
}

type MysqlUserRepository interface {
	Add(User) error
	GetById(id int) (User, error)
	GetUserCredentials(username string) (string, string, error)
}

type UserUsecase interface {
	Add(User) error
	GetById(id int) (User, error)
	GetUserCredentials(username string) (string, string, error)
}

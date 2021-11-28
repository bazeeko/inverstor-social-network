package domain

type User struct {
	ID               int      `json:"id"`
	Username         string   `json:"username"`
	Password         string   `json:"-"`
	Rating           float64  `json:"rating"`
	FavouriteTickers []string `json:"favourite_tickers"`
	CreatedAt        string   `json:"created_at"`
}

type MysqlUserRepository interface {
	Add(User) error
	GetById(id int) (User, error)
	GetByUsername(username string) (User, error)
	GetUserCredentials(username string) (string, string, error)
	AddUserToFavourites(userID int, favUserID int) error
	DeleteUserFromFavourites(userID int, favUserID int) error
	GetFavouriteUsers(userID int) ([]int, error)
	AddLikeToUser(userID int, likedUserID int) error
	DeleteLikeFromUser(userID int, likedUserID int) error
}

type UserUsecase interface {
	Add(User) error
	GetById(id int) (User, error)
	GetByUsername(username string) (User, error)
	GetUserCredentials(username string) (string, string, error)
	AddUserToFavourites(userID int, favUserID int) error
	DeleteUserFromFavourites(userID int, favUserID int) error
	GetFavouriteUsers(userID int) ([]User, error)
	AddLikeToUser(userID int, likedUserID int) error
	DeleteLikeFromUser(userID int, likedUserID int) error
}

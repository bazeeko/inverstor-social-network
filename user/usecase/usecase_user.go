package usercase

import (
	"fmt"

	"github.com/bazeeko/investor-social-network/domain"
)

type userUsecase struct {
	mysqlUserRepo domain.MysqlUserRepository
}

func NewUserUsecase(r domain.MysqlUserRepository) domain.UserUsecase {
	return &userUsecase{r}
}

func (uc *userUsecase) Add(u domain.User) error {
	_, err := uc.mysqlUserRepo.GetById(u.ID)
	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}

	err = uc.mysqlUserRepo.Add(u)
	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}

	return nil
}

func (uc *userUsecase) GetById(id int) (domain.User, error) {
	user, err := uc.mysqlUserRepo.GetById(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetById: %w", err)
	}

	// if user == (domain.User{}) {
	// 	return domain.User{}, fmt.Errorf("user not found")
	// }

	return user, nil
}

func (uc *userUsecase) GetByUsername(username string) (domain.User, error) {
	return uc.mysqlUserRepo.GetByUsername(username)
}

func (uc *userUsecase) GetUserCredentials(username string) (string, string, error) {
	return uc.mysqlUserRepo.GetUserCredentials(username)
}

func (uc *userUsecase) AddUserToFavourites(userID int, favUserID int) error {
	return uc.mysqlUserRepo.AddUserToFavourites(userID, favUserID)
}

func (uc *userUsecase) DeleteUserFromFavourites(userID int, favUserID int) error {
	return uc.mysqlUserRepo.DeleteUserFromFavourites(userID, favUserID)
}

func (uc *userUsecase) GetFavouriteUsers(userID int) ([]domain.User, error) {
	favUserIDs, err := uc.mysqlUserRepo.GetFavouriteUsers(userID)
	if err != nil {
		return []domain.User{}, fmt.Errorf("GetFavouriteUsers: %w", err)
	}

	favUsers := []domain.User{}

	for _, id := range favUserIDs {
		user, err := uc.mysqlUserRepo.GetById(id)
		if err != nil {
			return []domain.User{}, fmt.Errorf("GetFavouriteUsers: %w", err)
		}

		favUsers = append(favUsers, user)
	}

	return favUsers, nil
}

func (uc *userUsecase) AddLikeToUser(userID int, likedUserID int) error {
	return uc.mysqlUserRepo.AddLikeToUser(userID, likedUserID)
}

func (uc *userUsecase) DeleteLikeFromUser(userID int, likedUserID int) error {
	return uc.mysqlUserRepo.DeleteLikeFromUser(userID, likedUserID)
}

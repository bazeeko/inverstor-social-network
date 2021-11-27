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

	if user == (domain.User{}) {
		return domain.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (uc *userUsecase) GetByUsername(username string) (domain.User, error) {
	user, err := uc.mysqlUserRepo.GetByUsername(username)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetByUsername: %w", err)
	}

	if user == (domain.User{}) {
		return domain.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (uc *userUsecase) GetUserCredentials(username string) (string, string, error) {
	uname, pass, err := uc.mysqlUserRepo.GetUserCredentials(username)
	if err != nil {
		return "", "", fmt.Errorf("GetUserCredentials: %w", err)
	}

	return uname, pass, nil
}

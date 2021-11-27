package usecase

import (
	"fmt"

	"github.com/bazeeko/investor-social-network/domain"
)

type threadUsecase struct {
	mysqlThreadRepo domain.MysqlThreadRepository
}

func NewThreadUsecase(mtr domain.MysqlThreadRepository) domain.ThreadUsecase {
	return &threadUsecase{mtr}
}

func (uc *threadUsecase) CreateThread(userID int, t domain.Thread) error {
	err := uc.mysqlThreadRepo.CreateThread(userID, t)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

}

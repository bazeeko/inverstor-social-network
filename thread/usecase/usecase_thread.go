package usecase

import (
	"fmt"
	"log"

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
		return fmt.Errorf("CreateThread: %w", err)
	}

	return nil
}

func (uc *threadUsecase) GetThreadByID(threadID int) (domain.Thread, error) {
	t, err := uc.mysqlThreadRepo.GetThreadByID(threadID)
	if err != nil {
		return domain.Thread{}, fmt.Errorf("GetThreadByID: %w", err)
	}

	t.Comments, err = uc.mysqlThreadRepo.GetCommentsByThreadID(t.ID)
	if err != nil {
		log.Println(err)
	}

	for i := range t.Comments {
		t.Comments[i].SubComments, _ = uc.mysqlThreadRepo.GetSubCommentsByCommentID(t.Comments[i].ID)
		log.Println(err)
	}

	return t, nil
}

func (uc *threadUsecase) DeleteThreadByID(threadID int) error {
	err := uc.mysqlThreadRepo.DeleteThreadByID(threadID)
	if err != nil {
		return fmt.Errorf("DeleteThreadByID: %w", err)
	}

	return nil
}

func (uc *threadUsecase) GetUserThreads(userID int) ([]domain.Thread, error) {
	return uc.mysqlThreadRepo.GetUserThreads(userID)
}

func (uc *threadUsecase) GetThreadsByHashtag(hashtag string) ([]domain.Thread, error) {
	return uc.mysqlThreadRepo.GetThreadsByHashtag(hashtag)
}

func (uc *threadUsecase) CreateComment(c domain.Comment) error {
	if c.Body == "" {
		return fmt.Errorf("empty comment body")
	}
	return uc.mysqlThreadRepo.CreateComment(c)
}

func (uc *threadUsecase) GetCommentsByThreadID(threadID int) ([]domain.Comment, error) {
	return uc.mysqlThreadRepo.GetCommentsByThreadID(threadID)
}

func (uc *threadUsecase) CreateSubComment(sc domain.SubComment) error {
	if sc.Body == "" {
		return fmt.Errorf("empty comment body")
	}
	return uc.mysqlThreadRepo.CreateSubComment(sc)
}

func (uc *threadUsecase) GetSubCommentsByCommentID(commentID int) ([]domain.SubComment, error) {
	return uc.mysqlThreadRepo.GetSubCommentsByCommentID(commentID)
}

func (uc *threadUsecase) AddLikeToThread(userID int, threadID int) error {
	return uc.mysqlThreadRepo.AddLikeToThread(userID, threadID)
}

func (uc *threadUsecase) DeleteLikeFromThread(userID int, threadID int) error {
	return uc.mysqlThreadRepo.DeleteLikeFromThread(userID, threadID)
}

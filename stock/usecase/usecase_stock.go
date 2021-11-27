package UserUsecase

import (
	"fmt"

	"github.com/bazeeko/investor-social-network/domain"
)

type stockUsecase struct {
	mysqlStockRepo domain.MysqlStockRepository
}

func NewStockUsecase(msr domain.MysqlStockRepository) domain.StockUsecase {
	return &stockUsecase{msr}
}

func (uc *stockUsecase) GetStockBySymbol(symbol string) (domain.Stock, error) {
	s, err := uc.mysqlStockRepo.GetStockBySymbol(symbol)
	if err != nil {
		return domain.Stock{}, fmt.Errorf("GetStockBySymbol: %w", err)
	}

	if s == (domain.Stock{}) {
		return domain.Stock{}, fmt.Errorf("stock not found")
	}

	return s, nil
}

func (uc *stockUsecase) AddStockToFavourites(userID int, symbol string) error {
	err := uc.mysqlStockRepo.AddStockToFavourites(userID, symbol)
	if err != nil {
		return fmt.Errorf("AddStockToFavourites: %w", err)
	}

	return nil
}

func (uc *stockUsecase) DeleteStockFromFavourites(userID int, symbol string) error {
	err := uc.mysqlStockRepo.DeleteStockFromFavourites(userID, symbol)

	if err != nil {
		return fmt.Errorf("DeleteStockFromFavourites: %w", err)
	}

	return nil
}

func (uc *stockUsecase) GetFavouriteStocks(userID int) ([]string, error) {
	stocks, err := uc.mysqlStockRepo.GetFavouriteStocks(userID)
	if err != nil {
		return []string{}, err
	}

	return stocks, nil
}

// name, avatar, rating,

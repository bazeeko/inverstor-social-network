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

// name, avatar, rating,

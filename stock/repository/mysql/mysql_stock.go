package mysql

import (
	"database/sql"
	"fmt"

	"github.com/bazeeko/investor-social-network/domain"
)

type mysqlStockRepository struct {
	*sql.DB
}

func NewMysqlStockRepository(db *sql.DB) domain.MysqlStockRepository {
	return &mysqlStockRepository{db}
}

// func (r *mysqlStockRepository) AddStock(s domain.Stock) error {
// 	_, err := r.Exec(`INSERT stock (
// 		user_id, topic, body, image_url, created_at)
// 		VALUES (?, ?, ?, ?, ?)`,
// 		userID, t.Topic, t.Body, t.ImageURL, t.CreatedAt)

// 	if err != nil {
// 		return fmt.Errorf("CreateThread: %w", err)
// 	}

// 	return nil
// }

func (r *mysqlStockRepository) GetStockBySymbol(symbol string) (domain.Stock, error) {
	s := domain.Stock{}

	err := r.QueryRow(`SELECT
		symbol, name, info, image_url
		FROM stocks
		WHERE symbol=?`, symbol).
		Scan(&s.Symbol, &s.Name, &s.Info, &s.ImageURL)

	if err != nil {
		return domain.Stock{}, fmt.Errorf("GetStockBySymbol: %w", err)
	}

	return s, nil
}

func (r *mysqlStockRepository) AddStockToFavourites(userID int, symbol string) error {
	_, err := r.Exec(`INSERT favourite_stocks (
		user_id, stock_symbol)
		VALUES (?, ?)`,
		userID, symbol)

	if err != nil {
		return fmt.Errorf("AddStockToFavourites: %w", err)
	}

	return nil
}

func (r *mysqlStockRepository) DeleteStockFromFavourites(userID int, symbol string) error {
	_, err := r.Exec(`DELETE FROM favourite_stocks
		WHERE user_id=? AND stock_symbol=?`,
		userID, symbol)

	if err != nil {
		return fmt.Errorf("DeleteStockFromFavourites: %w", err)
	}

	return nil
}

func (r *mysqlStockRepository) GetFavouriteStocks(userID int) ([]string, error) {
	stocks := make([]string, 0)

	rows, err := r.Query(`SELECT stock_symbol FROM favourite_stocks
	WHERE user_id=?`, userID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var symbol string

		err := rows.Scan(&symbol)
		if err != nil {
			return []string{}, err
		}

		stocks = append(stocks, symbol)
	}

	if err != nil {
		return []string{}, fmt.Errorf("GetFavouriteStocks: %w", err)
	}

	return stocks, nil
}

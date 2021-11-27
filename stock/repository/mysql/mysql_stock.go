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

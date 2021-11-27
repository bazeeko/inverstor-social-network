package domain

type Stock struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Info     string `json:"info"`
	ImageURL string `json:"image_url"`
}

type MysqlStockRepository interface {
	GetStockBySymbol(symbol string) (Stock, error)
	AddStockToFavourites(userID int, symbol string) error
	DeleteStockFromFavourites(userID int, symbol string) error
	GetFavouriteStocks(userID int) ([]string, error)
}

type StockUsecase interface {
	GetStockBySymbol(symbol string) (Stock, error)
	AddStockToFavourites(userID int, symbol string) error
	DeleteStockFromFavourites(userID int, symbol string) error
	GetFavouriteStocks(userID int) ([]string, error)
}

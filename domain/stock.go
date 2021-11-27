package domain

type Stock struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Info     string `json:"info"`
	ImageURL string `json:"image_url"`
}

type MysqlStockRepository interface {
	GetStockBySymbol(symbol string) (Stock, error)
}

type StockUsecase interface {
	GetStockBySymbol(symbol string) (Stock, error)
}

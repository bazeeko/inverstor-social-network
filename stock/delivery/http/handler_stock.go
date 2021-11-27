package http

import (
	"log"
	"net/http"

	"github.com/bazeeko/investor-social-network/domain"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string `json:"message"`
}

type StockHandler struct {
	stockUsecase domain.StockUsecase
}

func NewStockHandler(e *echo.Echo, suc domain.StockUsecase) {
	handler := &StockHandler{suc}

	stockGroup := e.Group("/api/stock")

	stockGroup.GET("/:symbol", handler.GetBySymbol)
}

func (h *StockHandler) GetBySymbol(c echo.Context) error {
	symbol := c.Param("symbol")

	s, err := h.stockUsecase.GetStockBySymbol(symbol)
	if err != nil {
		log.Printf("GetBySymbol: %s\n", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal server error"})
	}

	return c.JSON(http.StatusOK, s)
}

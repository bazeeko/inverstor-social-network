package http

import (
	"context"
	"log"
	"net/http"

	"github.com/bazeeko/investor-social-network/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Message string `json:"message"`
}

type StockHandler struct {
	stockUsecase domain.StockUsecase
	userUsecase  domain.UserUsecase
}

type key struct {
}

func NewStockHandler(e *echo.Echo, suc domain.StockUsecase, uuc domain.UserUsecase) {
	handler := &StockHandler{suc, uuc}

	stockGroup := e.Group("/api/stock")

	stockGroup.Use(middleware.BasicAuth(func(s1, s2 string, c echo.Context) (bool, error) {
		username, password, err := handler.userUsecase.GetUserCredentials(s1)
		if err != nil {
			log.Printf("BasicAuth: %s\n", err)
			return false, nil
		}

		if username == s1 && password == s2 {
			ctx := context.WithValue(c.Request().Context(), key{}, username)
			c.SetRequest(c.Request().WithContext(ctx))

			return true, nil
		}

		return false, nil
	}))

	stockGroup.GET("/:symbol", handler.GetBySymbol)
	stockGroup.POST("/:symbol", handler.AddToFavourites)
	stockGroup.DELETE("/:symbol", handler.DeleteFromFavourites)
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

func (h *StockHandler) AddToFavourites(c echo.Context) error {
	symbol := c.Param("symbol")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("AddToFavourites: %s\n", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	err = h.stockUsecase.AddStockToFavourites(user.ID, symbol)
	if err != nil {
		log.Printf("AddToFavourites: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusOK)
}

func (h *StockHandler) DeleteFromFavourites(c echo.Context) error {
	symbol := c.Param("symbol")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("DeleteFromFavourites: %s\n", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	err = h.stockUsecase.DeleteStockFromFavourites(user.ID, symbol)
	if err != nil {
		log.Printf("DeleteFromFavourites: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusOK)
}

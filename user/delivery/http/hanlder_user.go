package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/bazeeko/investor-social-network/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Message string `json:"message"`
}

type UserHandler struct {
	userUsecase  domain.UserUsecase
	stockUsecase domain.StockUsecase
}

type key struct {
}

func NewUserHandler(e *echo.Echo, uuc domain.UserUsecase, suc domain.StockUsecase) {
	handler := &UserHandler{uuc, suc}

	userGroup := e.Group("/api/user")

	userGroup.Use(middleware.BasicAuth(func(s1, s2 string, c echo.Context) (bool, error) {
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

	userGroup.GET("/:id", handler.GetUser)
	userGroup.GET("/:id/favourite/tickers", handler.GetFavStocks)
	userGroup.GET("/:id/favourite/users", handler.GetFavUsers)
	userGroup.POST("/:id/addToFavourites", handler.GetFavUsers)
	userGroup.DELETE("/:id/deleteFromFavourites", handler.GetFavUsers)
	userGroup.POST("/:id/like", handler.LikeUser)
	userGroup.DELETE("/:id/like", handler.DeleteLikeUser)
	// userGroup.POST("/:id", handler.AddUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	}

	user, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("GetUser: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"User Not Found"})
	}

	user.FavouriteTickers, err = h.stockUsecase.GetFavouriteStocks(user.ID)
	if err != nil {
		log.Printf("GetUser: %s\n", err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetFavStocks(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetFavStocks: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	}

	// username, ok := c.Request().Context().Value(key{}).(string)
	// if !ok {
	// 	return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	// }

	user, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("GetFavStocks: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"User Not Found"})
	}

	stocks, err := h.stockUsecase.GetFavouriteStocks(user.ID)
	if err != nil {
		log.Printf("GetFavStocks: %s\n", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	return c.JSON(http.StatusOK, stocks)
}

func (h *UserHandler) GetFavUsers(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetFavUsers: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	}

	// username, ok := c.Request().Context().Value(key{}).(string)
	// if !ok {
	// 	return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	// }

	user, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("GetFavUsers: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"User Not Found"})
	}

	favUsers, err := h.userUsecase.GetFavouriteUsers(user.ID)
	if err != nil {
		log.Printf("GetFavUsers: %s", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	return c.JSON(http.StatusOK, favUsers)
}

func (h *UserHandler) LikeUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("LikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	}

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user1, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("LikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	user2, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("LikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	err = h.userUsecase.AddLikeToUser(user1.ID, user2.ID)
	if err != nil {
		log.Printf("LikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) DeleteLikeUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("DeleteLikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	}

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user1, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("DeleteLikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	user2, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("DeleteLikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	err = h.userUsecase.DeleteLikeFromUser(user1.ID, user2.ID)
	if err != nil {
		log.Printf("DeleteLikeUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusOK)
}

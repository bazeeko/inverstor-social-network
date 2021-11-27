package http

import (
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
	userUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, uuc domain.UserUsecase) {
	handler := &UserHandler{uuc}

	userGroup := e.Group("/api/user")

	userGroup.Use(middleware.BasicAuth(func(s1, s2 string, c echo.Context) (bool, error) {
		username, password, err := handler.userUsecase.GetUserCredentials(s1)
		if err != nil {
			log.Printf("BasicAuth: %s\n", err)
			return false, nil
		}

		if username == s1 && password == s2 {
			return true, nil
		}

		return false, nil
	}))

	userGroup.GET("/:id", handler.GetUser)
	// userGroup.POST("/:id", handler.AddUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"invalid user id"})
	}

	user, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("GetUser: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"some error"})
	}

	return c.JSON(http.StatusOK, user)
}

// func (h *UserHandler) AddUser(c echo.Context) error {
// 	return nil
// }

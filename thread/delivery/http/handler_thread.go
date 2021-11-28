package http

import (
	"context"
	"encoding/json"
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

type ThreadHandler struct {
	threadUsecase domain.ThreadUsecase
	userUsecase   domain.UserUsecase
}

type key struct {
}

func NewThreadHandler(e echo.Echo, tuc domain.ThreadUsecase, uuc domain.UserUsecase) {
	handler := &ThreadHandler{tuc, uuc}

	threadGroup := e.Group("/api/thread")

	threadGroup.Use(middleware.BasicAuth(func(s1, s2 string, c echo.Context) (bool, error) {
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

	threadGroup.POST("/", handler.PostThread)
	threadGroup.GET("/:id", handler.GetByID)
	threadGroup.DELETE("/:id", handler.DeleteByID)
	threadGroup.GET("/user/:id", handler.GetThreadsByUserID)
	threadGroup.GET("/hashtag/:hashtag", handler.GetByHashtag)

	threadGroup.POST("/:id/comment/", handler.PostComment)
	threadGroup.POST("/:thread_id/comment/:comment_id/subcomment/", handler.PostSubComment)

	threadGroup.POST("/:id/like", handler.LikeThread)
	threadGroup.DELETE("/:id/like", handler.DeleteLikeThread)

	// userGroup.GET("/:id", handler.GetUser)
	// userGroup.GET("/:id/favourite/tickers", handler.GetFavStocks)
	// userGroup.GET("/:id/favourite/users", handler.GetFavUsers)
	// userGroup.POST("/:id/addToFavourites", handler.GetFavUsers)
	// userGroup.DELETE("/:id/deleteFromFavourites", handler.GetFavUsers)
	// userGroup.POST("/:id/like", handler.LikeUser)
	// userGroup.DELETE("/:id/like", handler.DeleteLikeUser)
}

func (h *ThreadHandler) PostThread(c echo.Context) error {
	// id, err := strconv.Atoi(c.Param("id"))

	// if err != nil {
	// 	log.Printf("PostThread: %s\n", err)
	// 	return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	// }

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("PostThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	t := domain.Thread{}

	err = json.NewDecoder(c.Request().Body).Decode(&t)
	if err != nil {
		log.Printf("PostThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	err = h.threadUsecase.CreateThread(user.ID, t)
	if err != nil {
		log.Printf("PostThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *ThreadHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetByID: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Thread Id"})
	}

	t, err := h.threadUsecase.GetThreadByID(id)
	if err != nil {
		log.Printf("GetByID: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"Thread Not Found"})
	}

	return c.JSON(http.StatusOK, t)
}

func (h *ThreadHandler) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("DeleteByID: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Thread Id"})
	}

	err = h.threadUsecase.DeleteThreadByID(id)
	if err != nil {
		log.Printf("DeleteByID: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"Thread Not Found"})
	}

	return c.NoContent(http.StatusOK)
}

func (h *ThreadHandler) GetThreadsByUserID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetThreadsByUserID: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid User Id"})
	}

	_, err = h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("GetThreadsByUserID: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"User Not Found"})
	}

	ts, err := h.threadUsecase.GetUserThreads(id)
	if err != nil {
		log.Printf("GetThreadsByUserID: %s\n", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *ThreadHandler) GetByHashtag(c echo.Context) error {
	hashtag := c.Param("hashtag")

	ts, err := h.threadUsecase.GetThreadsByHashtag(hashtag)
	if err != nil {
		log.Printf("GetByHashtag: %s\n", err)
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *ThreadHandler) PostComment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("PostComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Thread Id"})
	}

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("PostComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	_, err = h.threadUsecase.GetThreadByID(id)
	if err != nil {
		log.Printf("PostComment: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"Thread Not Found"})
	}

	comment := domain.Comment{}

	err = json.NewDecoder(c.Request().Body).Decode(&comment)
	if err != nil {
		log.Printf("PostComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	comment.ThreadID = id
	comment.UserID = user.ID

	err = h.threadUsecase.CreateComment(comment)
	if err != nil {
		log.Printf("PostComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *ThreadHandler) PostSubComment(c echo.Context) error {
	threadID, err := strconv.Atoi(c.Param("thread_id"))

	if err != nil {
		log.Printf("PostSubComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Thread Id"})
	}

	commentID, err := strconv.Atoi(c.Param("comment_id"))

	if err != nil {
		log.Printf("PostSubComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Comment Id"})
	}

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("PostSubComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	_, err = h.threadUsecase.GetThreadByID(threadID)
	if err != nil {
		log.Printf("PostSubComment: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"Thread Not Found"})
	}

	subComment := domain.SubComment{}

	err = json.NewDecoder(c.Request().Body).Decode(&subComment)
	if err != nil {
		log.Printf("PostSubComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	subComment.ThreadID = threadID
	subComment.CommentID = commentID
	subComment.UserID = user.ID

	err = h.threadUsecase.CreateSubComment(subComment)
	if err != nil {
		log.Printf("PostSubComment: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *ThreadHandler) LikeThread(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("LikeThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Thread Id"})
	}

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("LikeThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	_, err = h.threadUsecase.GetThreadByID(id)
	if err != nil {
		log.Printf("LikeThread: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"Thread Not Found"})
	}

	err = h.threadUsecase.AddLikeToThread(user.ID, id)
	if err != nil {
		log.Printf("LikeThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusOK)
}

func (h *ThreadHandler) DeleteLikeThread(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("DeleteLikeThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Invalid Thread Id"})
	}

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, Response{"Internal Server Error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("DeleteLikeThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	_, err = h.threadUsecase.GetThreadByID(id)
	if err != nil {
		log.Printf("DeleteLikeThread: %s\n", err)
		return c.JSON(http.StatusNotFound, Response{"Thread Not Found"})
	}

	err = h.threadUsecase.DeleteLikeFromThread(user.ID, id)
	if err != nil {
		log.Printf("DeleteLikeThread: %s\n", err)
		return c.JSON(http.StatusBadRequest, Response{"Bad Request"})
	}

	return c.NoContent(http.StatusOK)
}
